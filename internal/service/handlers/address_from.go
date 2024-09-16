package handlers

import (
	"github.com/kish1n/usdt_listening/internal/service/errors/apierrors"
	"github.com/kish1n/usdt_listening/internal/service/helpers"
	"gitlab.com/distributed_lab/ape"
	"net/http"
)

func SortBySender(w http.ResponseWriter, r *http.Request) {
	logger := Log(r)
	db := DB(r)
	address, err := helpers.GetAddress(r, "from_address")
	res, err := db.Link().SortByParameter(address, "from_address")

	if res == nil {
		apierrors.ErrorConstructor(w, *logger, err, "404 not found", "404", "Not Found", "Not found transaction from this address")
		return
	}

	if err != nil {
		apierrors.ErrorConstructor(w, *logger, err, "Server error", "500", "Server error 500", "Unpredictable behavior")
		return
	}

	logger.Infof("res: %s", res)
	ape.Render(w, res)
	return
}
