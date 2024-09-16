package handlers

import (
	"context"
	"net/http"

	"github.com/kish1n/usdt_listening/internal/data"

	"gitlab.com/distributed_lab/logan/v3"
)

type ctxKey int

const (
	logCtxKey ctxKey = iota
	dbCtxKey
)

func CtxLog(entry *logan.Entry) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, logCtxKey, entry)
	}
}

func TransactionQ(r *http.Request) data.TransactionQ {
	return r.Context().Value(dbCtxKey).(data.TransactionQ).New()
}

func CtxTransactionQ(q data.TransactionQ) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, dbCtxKey, q)
	}
}
func Log(r *http.Request) *logan.Entry {
	return r.Context().Value(logCtxKey).(*logan.Entry)
}
