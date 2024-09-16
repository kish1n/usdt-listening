package pg

import (
	"database/sql"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/kish1n/usdt_listening/internal/data"
	"gitlab.com/distributed_lab/kit/pgdb"
)

const TransactionTable = "transfers"

type transaction struct {
	db       *pgdb.DB
	selector sq.SelectBuilder
	updater  sq.UpdateBuilder
	counter  sq.SelectBuilder
}

func NewTransaction(db *pgdb.DB) data.TransactionQ {
	return &transaction{
		db:       db,
		selector: sq.Select("*").From(TransactionTable),
		updater:  sq.Update(TransactionTable),
		counter:  sq.Select("COUNT(*) as count").From(TransactionTable),
	}
}

func (q *transaction) New() data.TransactionQ {
	return NewTransaction(q.db)
}

func (q *transaction) Get() (*data.Transaction, error) {
	var res data.Transaction

	if err := q.db.Get(&res, q.selector); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get balance: %w", err)
	}

	return &res, nil
}

func (q *transaction) Select() ([]data.Transaction, error) {
	var res []data.Transaction

	if err := q.db.Select(&res, q.selector); err != nil {
		return nil, fmt.Errorf("select balances: %w", err)
	}

	return res, nil
}

func (q *transaction) Insert(trn data.Transaction) error {
	stmt := sq.Insert(TransactionTable).SetMap(map[string]interface{}{
		"id":           trn.Id,
		"from_address": trn.FromAddress,
		"to_address":   trn.ToAddress,
		"value":        trn.Value,
		"created_at":   trn.CreatedAt,
	})

	if err := q.db.Exec(stmt); err != nil {
		return fmt.Errorf("insert %s %+v: %w", TransactionTable, trn, err)
	}

	return nil
}

func (q *transaction) FilterByFromAddress(address string) data.TransactionQ {
	return q.applyCondition(sq.Eq{"from_address": address})
}

func (q *transaction) FilterByToAddress(address string) data.TransactionQ {
	return q.applyCondition(sq.Eq{"to_address": address})
}

func (q *transaction) FilterByAddress(address string) data.TransactionQ {
	res := q.applyCondition(sq.Or{
		sq.Eq{"to_address": address},
		sq.Eq{"from_address": address},
	})
	return res
}

func (q *transaction) PageBySide(page *pgdb.OffsetPageParams, column string) data.TransactionQ {
	q.selector = page.ApplyTo(q.selector, column)
	return q
}

func (q *transaction) Page(page *pgdb.OffsetPageParams) data.TransactionQ {
	q.selector = page.ApplyTo(q.selector, "created_at")
	return q
}

func (q *transaction) Count() (int64, error) {
	res := struct {
		Count int64 `db:"count"`
	}{}

	if err := q.db.Get(&res, q.counter); err != nil {
		return 0, fmt.Errorf("count transaction: %w", err)
	}

	return res.Count, nil
}

func (q *transaction) applyCondition(cond sq.Sqlizer) data.TransactionQ {
	q.selector = q.selector.Where(cond)
	q.updater = q.updater.Where(cond)
	q.counter = q.counter.Where(cond)
	return q
}
