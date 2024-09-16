package service

import (
	"net"

	"github.com/kish1n/usdt_listening/internal/config"
	"gitlab.com/distributed_lab/kit/copus/types"
	"gitlab.com/distributed_lab/logan/v3"
)

type service struct {
	log      *logan.Entry
	copus    types.Copus
	listener net.Listener
}

//func (s *service) run(ctx context.Context, cfg config.Config) error {
//	s.log.Info("Service started")
//	r, err := s.router(ctx, cfg)
//
//	if err != nil {
//		s.log.Error(err.Error())
//		return err
//	}
//
//	if err := s.copus.RegisterChi(r); err != nil {
//		return errors.Wrap(err, "cop failed")
//	}
//
//	return http.Serve(s.listener, r)
//}

func newService(cfg config.Config) *service {
	return &service{
		log:      cfg.Log(),
		copus:    cfg.Copus(),
		listener: cfg.Listener(),
	}
}

//func Run(ctx context.Context, cfg config.Config) {
//	if err := newService(cfg).run(ctx, cfg); err != nil {
//		panic(err)
//	}
//	log := cfg.Log()
//
//	signalChan := make(chan os.Signal, 1)
//	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
//
//	select {
//	case <-signalChan:
//		log.Info("Received shutdown signal, shutting down...")
//	case <-ctx.Done():
//		log.Info("Context cancelled, shutting down...")
//	}
//
//	log.Info("Process terminated gracefully")
//}
