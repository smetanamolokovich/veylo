package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/smetanamolokovich/veylo/internal/interface/http/handler"
	authmiddleware "github.com/smetanamolokovich/veylo/internal/interface/http/middleware"
	"github.com/smetanamolokovich/veylo/pkg/jwt"
)

func NewRouter(inspectionHandler *handler.InspectionHandler, jwtManager *jwt.Manager) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/inspections", func(r chi.Router) {
			r.Use(authmiddleware.Auth(jwtManager))
			r.Post("/", inspectionHandler.Create)
		})
	})

	return r
}
