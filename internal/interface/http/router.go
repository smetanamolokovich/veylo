package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/smetanamolokovich/veylo/internal/interface/http/handler"
	authmiddleware "github.com/smetanamolokovich/veylo/internal/interface/http/middleware"
	"github.com/smetanamolokovich/veylo/pkg/jwt"
)

func NewRouter(inspectionHandler *handler.InspectionHandler, authHandler *handler.AuthHandler, assetHandler *handler.AssetHandler, findingHandler *handler.FindingHandler, workflowHandler *handler.WorkflowHandler, orgHandler *handler.OrganizationHandler, invitationHandler *handler.InvitationHandler, jwtManager *jwt.Manager) *chi.Mux {
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
			r.Get("/{id}/report", inspectionHandler.GetReport)
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
		r.Route("/workflow", func(r chi.Router) {
			r.Use(authmiddleware.Auth(jwtManager))
			r.Post("/", workflowHandler.Create)
			r.Get("/", workflowHandler.Get)
			r.Post("/statuses", workflowHandler.AddStatus)
			r.Post("/transitions", workflowHandler.AddTransition)
		})
	})

	r.Route("/api/v1/organizations", func(r chi.Router) {
		r.Use(authmiddleware.Auth(jwtManager))
		r.Post("/", orgHandler.Create)
		r.Get("/me", orgHandler.GetMe)
		r.Post("/me/onboarding", orgHandler.CompleteOnboarding)
		r.Post("/me/invitations", invitationHandler.Create)
	})

	r.Route("/api/auth", func(r chi.Router) {
		r.Post("/signup", authHandler.Signup)
		r.Post("/register", authHandler.Register)
		r.Post("/login", authHandler.Login)
		r.Post("/refresh", authHandler.Refresh)
		r.Get("/invite/{token}", invitationHandler.GetByToken)
		r.Post("/invite/{token}/accept", invitationHandler.Accept)
	})

	return r
}
