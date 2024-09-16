package service

import (
	"github.com/go-chi/chi"
	"github.com/kish1n/usdt_listening/internal/config"
	"github.com/kish1n/usdt_listening/internal/data/pg"
	"github.com/kish1n/usdt_listening/internal/service/handlers"
	"gitlab.com/distributed_lab/ape"
	"net/http"
)

func (s *service) router(cfg config.Config) (chi.Router, error) {
	r := chi.NewRouter()
	logger := cfg.Log()

	r.Use(
		ape.RecoverMiddleware(s.log),
		ape.LoganMiddleware(s.log),
		ape.CtxMiddleware(
			handlers.CtxLog(s.log),
			handlers.CtxDB(pg.NewMasterQ(cfg.DB())),
		),
	)

	r.Route("/", func(r chi.Router) {
		r.Get("/listen", handlers.ListenForTransfers)
		r.Get("/from/{from_address}", handlers.SortBySender)
		r.Get("/to/{to_address}", handlers.SortByRecipient)
		r.Get("/by/{address}", handlers.SortByAddress)
	})

	logger.Info("Starting server on :8080")
	err := http.ListenAndServe(":8080", r)

	if err != nil {
		logger.Fatalf("Failed to start server: %v", err)
		return r, err
	}

	return r, nil
}
