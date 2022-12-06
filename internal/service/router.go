package service

import (
	"github.com/mhrynenko/jwt_service/internal/data/postgres"
	"github.com/mhrynenko/jwt_service/internal/service/handlers"

	"github.com/go-chi/chi"
	"gitlab.com/distributed_lab/ape"
)

func (s *service) router() chi.Router {
	r := chi.NewRouter()

	r.Use(
		ape.RecoverMiddleware(s.log),
		ape.LoganMiddleware(s.log),
		ape.CtxMiddleware(
			handlers.CtxLog(s.log),
			handlers.CtxUsersQ(postgres.NewUsersQ(s.cfg.DB())),
			handlers.CtxRefreshTokensQ(postgres.NewRefreshTokensQ(s.cfg.DB())),
		),
	)
	r.Route("/login", func(r chi.Router) {
		r.Post("/", handlers.Login)
	})

	r.Route("/refresh", func(r chi.Router) {
		r.Post("/", handlers.Refresh)
	})

	r.Route("/logout", func(r chi.Router) {
		r.Post("/", handlers.Logout)
	})

	return r
}
