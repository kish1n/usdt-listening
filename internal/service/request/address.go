package request

import (
	"net/http"

	"github.com/kish1n/usdt_listening/internal/service/page"
	"gitlab.com/distributed_lab/urlval/v4"
)

type TransactionList struct {
	page.OffsetParams
	Count bool `url:"count"`
}

func GetAddress(r *http.Request) (req TransactionList, err error) {
	if err = urlval.Decode(r.URL.Query(), &req); err != nil {
		err = newDecodeError("query", err)
		return
	}

	return req, req.Validate()
}
