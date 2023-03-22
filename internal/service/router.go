package service

import (
	"context"

	"github.com/go-chi/chi"
	"gitlab.com/distributed_lab/acs/auth/internal/data/postgres"
	"gitlab.com/distributed_lab/acs/auth/internal/service/handlers"
	"gitlab.com/distributed_lab/acs/auth/internal/service/receiver"
	"gitlab.com/distributed_lab/acs/auth/internal/service/sender"
	"gitlab.com/distributed_lab/ape"
)

func (s *service) router() chi.Router {
	r := chi.NewRouter()
	ctx := context.Background()

	s.startSender(ctx)
	s.startReceiver(ctx)

	r.Use(
		ape.RecoverMiddleware(s.log),
		ape.LoganMiddleware(s.log),
		ape.CtxMiddleware(
			handlers.CtxLog(s.log),
			handlers.CtxJwtParams(*s.cfg.JwtParams()),
			handlers.CtxUsersQ(postgres.NewUsersQ(s.cfg.DB())),
			handlers.CtxModulesQ(postgres.NewModulesQ(s.cfg.DB())),
			handlers.CtxPermissionsQ(postgres.NewPermissionsQ(s.cfg.DB())),
			handlers.CtxRefreshTokensQ(postgres.NewRefreshTokensQ(s.cfg.DB())),
		),
	)

	r.Route("/integrations/auth", func(r chi.Router) {
		r.Route("/login", func(r chi.Router) {
			r.Post("/", handlers.Login)
		})

		r.Route("/refresh", func(r chi.Router) {
			r.Post("/", handlers.Refresh)
		})

		r.Route("/logout", func(r chi.Router) {
			r.Post("/", handlers.Logout)
		})

		r.Route("/validate", func(r chi.Router) {
			r.Post("/", handlers.Validate)
		})
	})

	return r
}

func (s *service) startReceiver(ctx context.Context) {
	s.log.Info("Starting receiver")
	receiver.Run(ctx, s.cfg)
}

func (s *service) startSender(ctx context.Context) {
	s.log.Info("Starting sender")
	sender.Run(ctx, s.cfg)
}
