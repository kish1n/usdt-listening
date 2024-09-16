package service

import (
	"context"

	"github.com/go-chi/chi"
	"github.com/kish1n/usdt_listening/internal/config"
	"github.com/kish1n/usdt_listening/internal/service/handlers"
	"gitlab.com/distributed_lab/ape"
)

func Router(ctx context.Context, cfg config.Config) {
	r := chi.NewRouter()

	r.Use(
		ape.RecoverMiddleware(cfg.Log()),
		ape.LoganMiddleware(cfg.Log()),
		ape.CtxMiddleware(
			handlers.CtxLog(cfg.Log()),
		),
		handlers.DBCloneMiddleware(cfg.DB()),
	)

	r.Route("/transactions", func(r chi.Router) {
		r.Get("/from/{address}", handlers.SortBySender)
		r.Get("/to/{address}", handlers.SortByOrder)
		r.Get("/by/{address}", handlers.SortByAddress)
	})

	cfg.Log().Info("Service started")
	ape.Serve(ctx, r, cfg, ape.ServeOpts{})
}
