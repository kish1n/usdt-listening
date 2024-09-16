package data

import (
	"time"
)

type TransactionQ interface {
	Insert(trn TransactionData) error
	FilterByFromAddress(address string) TransactionQ
	FilterByToAddress(address string) TransactionQ
	FilterByAddress(address string) TransactionQ
}

type TransactionData struct {
	FromAddress string    `db:"from_address" json:"id"`
	ToAddress   string    `db:"to_address"   json:"from_address"`
	Value       int64     `db:"value"        json:"to_address"`
	Id          string    `db:"id"           json:"value"`
	CreatedAt   time.Time `db:"created_at"    json:"created_at"`
}
