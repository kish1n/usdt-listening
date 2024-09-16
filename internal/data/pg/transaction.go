package pg

import (
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/kish1n/usdt_listening/internal/data"
	"gitlab.com/distributed_lab/kit/pgdb"
)

const TransactionTable = "transfers"

type TransactionQ struct {
	db       *pgdb.DB
	selector sq.SelectBuilder
	updater  sq.UpdateBuilder
	counter  sq.SelectBuilder
}

func newTransactionQ(db *pgdb.DB) data.TransactionQ {
	return &TransactionQ{
		db:       db,
		selector: sq.Select("*").From(TransactionTable),
		updater:  sq.Update(TransactionTable),
		counter:  sq.Select("COUNT(*) as count").From(TransactionTable),
	}
}

func (q *TransactionQ) Insert(trn data.TransactionData) error {
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

func (q *TransactionQ) FilterByFromAddress(address string) data.TransactionQ {
	return q.applyCondition(sq.Eq{"from_address": address})
}

func (q *TransactionQ) FilterByToAddress(address string) data.TransactionQ {
	return q.applyCondition(sq.Eq{"to_address": address})
}

func (q *TransactionQ) FilterByAddress(address string) data.TransactionQ {
	res := q.applyCondition(sq.Or{
		sq.Eq{"to_address": address},
		sq.Eq{"from_address": address},
	})
	return res
}

func (q *TransactionQ) applyCondition(cond sq.Sqlizer) data.TransactionQ {
	q.selector = q.selector.Where(cond)
	q.updater = q.updater.Where(cond)
	q.counter = q.counter.Where(cond)
	return q
}
