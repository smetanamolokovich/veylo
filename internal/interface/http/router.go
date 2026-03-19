package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/smetanamolokovich/veylo/internal/interface/http/handler"
	authmiddleware "github.com/smetanamolokovich/veylo/internal/interface/http/middleware"
	"github.com/smetanamolokovich/veylo/pkg/jwt"
)

func NewRouter(inspectionHandler *handler.InspectionHandler, authHandler *handler.AuthHandler, assetHandler *handler.AssetHandler, findingHandler *handler.FindingHandler, jwtManager *jwt.Manager) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/inspections", func(r chi.Router) {
			r.Use(authmiddleware.Auth(jwtManager))
			r.Post("/", inspectionHandler.Create)
			r.Get("/", inspectionHandler.List)
			r.Get("/{id}", inspectionHandler.Get)
			r.Post("/{id}/transitions", inspectionHandler.Transition)
			r.Route("/{inspectionID}/findings", func(r chi.Router) {
				r.Post("/", findingHandler.Create)
				r.Get("/", findingHandler.List)
				r.Put("/{id}/assessment", findingHandler.Assess)
			})
		})
		r.Route("/assets", func(r chi.Router) {
			r.Use(authmiddleware.Auth(jwtManager))
			r.Post("/vehicles", assetHandler.CreateVehicle)
			r.Get("/{id}", assetHandler.Get)
		})
	})

	r.Route("/api/auth", func(r chi.Router) {
		r.Post("/register", authHandler.Register)
		r.Post("/login", authHandler.Login)
		r.Post("/refresh", authHandler.Refresh)
	})

	return r
}
