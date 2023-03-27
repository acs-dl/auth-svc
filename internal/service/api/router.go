package api

import (
	"github.com/go-chi/chi"
	"gitlab.com/distributed_lab/acs/auth/internal/data/postgres"
	"gitlab.com/distributed_lab/acs/auth/internal/service/api/handlers"
	"gitlab.com/distributed_lab/ape"
)

func (s *service) router() chi.Router {
	r := chi.NewRouter()

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
