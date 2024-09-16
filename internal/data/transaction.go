package data

import (
	"time"

	"gitlab.com/distributed_lab/kit/pgdb"
)

const (
	Sender = "from_address"
	Order  = "to_address"
)

type TransactionQ interface {
	New() TransactionQ
	Get() (*Transaction, error)
	Select() ([]Transaction, error)
	Insert(trn Transaction) error
	FilterByFromAddress(address string) TransactionQ
	FilterByToAddress(address string) TransactionQ
	FilterByAddress(address string) TransactionQ
	Page(page *pgdb.OffsetPageParams) TransactionQ
	PageBySide(page *pgdb.OffsetPageParams, column string) TransactionQ
	Count() (int64, error)
}

type Transaction struct {
	FromAddress string    `db:"from_address" json:"id"`
	ToAddress   string    `db:"to_address"   json:"from_address"`
	Value       int64     `db:"value"        json:"to_address"`
	Id          string    `db:"id"           json:"value"`
	CreatedAt   time.Time `db:"created_at"    json:"created_at"`
}
