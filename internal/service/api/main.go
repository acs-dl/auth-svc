package api

import (
	"context"
	"gitlab.com/distributed_lab/logan/v3"
	"net"
	"net/http"

	"github.com/acs-dl/auth-svc/internal/config"

	"gitlab.com/distributed_lab/kit/copus/types"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type service struct {
	log      *logan.Entry
	copus    types.Copus
	listener net.Listener
	cfg      config.Config
}

func (s *service) run() error {
	s.log.Info("Service started")
	r := s.router()

	if err := s.copus.RegisterChi(r); err != nil {
		return errors.Wrap(err, "cop failed")
	}

	return http.Serve(s.listener, r)
}

func newService(cfg config.Config) *service {
	return &service{
		log:      cfg.Log(),
		copus:    cfg.Copus(),
		listener: cfg.Listener(),
		cfg:      cfg,
	}
}

func Run(_ context.Context, cfg config.Config) {
	if err := newService(cfg).run(); err != nil {
		panic(err)
	}
}
