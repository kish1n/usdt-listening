package handlers

import (
	"strings"

	"github.com/go-chi/chi"
	"github.com/kish1n/usdt_listening/internal/data"
	"github.com/kish1n/usdt_listening/internal/service/request"
	"github.com/kish1n/usdt_listening/resources"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"

	"net/http"
)

func SortBySender(w http.ResponseWriter, r *http.Request) {
	address := strings.ToLower(chi.URLParam(r, "address"))

	req, err := request.GetAddress(r)
	if err != nil {
		Log(r).WithError(err).Error("error getting address")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	transactions, err := TransactionQ(r).Page(&req.OffsetPageParams).Select()
	if err != nil {
		Log(r).WithError(err).Error("Error getting transaction")
		ape.RenderErr(w, problems.InternalError())
		return
	}
	if transactions == nil {
		Log(r).Error("Transactions by this address:%s not found", address)
		ape.RenderErr(w, problems.NotFound())
		return
	}

	resp, err := NewTransactionResponseList(transactions)
	if err != nil {
		Log(r).WithError(err).Error("Error creating transaction list")
		ape.RenderErr(w, problems.InternalError())
		return
	}
	TrxsCount, err := TransactionQ(r).FilterByFromAddress(address).Count()
	if err != nil {
		Log(r).WithError(err).Error("Error getting transaction count")
		ape.RenderErr(w, problems.InternalError())
		return
	}
	resp.Links = req.GetLinks(r, uint64(TrxsCount))
	if req.Count {
		_ = resp.PutMeta(struct {
			TransactionCount int64 `json:"transaction_count"`
		}{TrxsCount})
	}
	ape.Render(w, resp)
}

func NewTransactionResponseList(trxs []data.Transaction) (resources.TransactionListResponse, error) {
	list := make([]resources.Transaction, len(trxs))
	for i, trx := range trxs {
		list[i] = resources.Transaction{
			Key: resources.Key{
				ID:   trx.Id,
				Type: resources.TRANSACTION,
			},
			Attributes: resources.TransactionAttributes{
				CreatedAt:   trx.CreatedAt,
				FromAddress: trx.FromAddress,
				ToAddress:   trx.ToAddress,
				Value:       int(trx.Value),
			},
		}
	}
	return resources.TransactionListResponse{Data: list}, nil
}
