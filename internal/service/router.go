package service

import (
	"github.com/go-chi/chi"
	"gitlab.com/distributed_lab/Auth/internal/data/postgres"
	"gitlab.com/distributed_lab/Auth/internal/service/handlers"
	middleware "gitlab.com/distributed_lab/Auth/internal/service/middlewares"
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
			handlers.CtxAmountsQ(postgres.NewAmountsQ(s.cfg.DB())),
			handlers.CtxModulesQ(postgres.NewModulesQ(s.cfg.DB())),
			handlers.CtxModulesUsersQ(postgres.NewModulesUsersQ(s.cfg.DB())),
			handlers.CtxRefreshTokensQ(postgres.NewRefreshTokensQ(s.cfg.DB())),
		),
	)
	r.Route("/login", func(r chi.Router) {
		r.Post("/", handlers.Login)
	})

	r.Route("/refresh", func(r chi.Router) {
		r.Use(middleware.Jwt())
		r.Post("/", handlers.Refresh)
	})

	r.Route("/logout", func(r chi.Router) {
		r.Post("/", handlers.Logout)
	})

	r.Route("/module", func(r chi.Router) {
		r.Post("/", handlers.AddModule)
		r.Route("/{name}", func(r chi.Router) {
			r.Get("/", handlers.GetModule)
			r.Delete("/", handlers.DeleteModule)
		})
	})

	r.Route("/permission", func(r chi.Router) {
		r.Post("/", handlers.AddModuleUser)
		r.Delete("/", handlers.DeleteModuleUser)
	})

	return r
}
