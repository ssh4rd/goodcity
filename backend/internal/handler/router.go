package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"goodcity/backend/internal/domain"
	"goodcity/backend/internal/service"
)

func NewRouter(
	authSvc *service.AuthService,
	practiceSvc *service.PracticeService,
	ratingSvc *service.RatingService,
) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
	}))

	authHandler := NewAuthHandler(authSvc)
	practiceHandler := NewPracticeHandler(practiceSvc)
	moderationHandler := NewModerationHandler(practiceSvc)
	ratingHandler := NewRatingHandler(ratingSvc)

	auth := AuthMiddleware(authSvc)
	moderatorOnly := RequireRole(domain.RoleModerator)

	r.Route("/api", func(r chi.Router) {
		// Auth
		r.Post("/auth/register", authHandler.Register)
		r.Post("/auth/login", authHandler.Login)

		// Practices (public)
		r.Get("/practices", practiceHandler.List)
		r.Get("/practices/{id}", practiceHandler.Get)
		r.Get("/practices/{id}/comments", practiceHandler.ListComments)
		r.Get("/practices/{id}/rating", ratingHandler.GetStats)

		// Practices (authenticated)
		r.Group(func(r chi.Router) {
			r.Use(auth)
			r.Post("/practices", practiceHandler.Create)
			r.Post("/practices/{id}/rate", ratingHandler.Rate)
			r.Get("/practices/{id}/rate", ratingHandler.GetMyRating)
			r.Post("/practices/{id}/comments", practiceHandler.AddComment)
		})

		// Practices (moderator)
		r.Group(func(r chi.Router) {
			r.Use(auth, moderatorOnly)
			r.Put("/practices/{id}", practiceHandler.Update)
			r.Delete("/practices/{id}", practiceHandler.Delete)
		})

		// Moderation
		r.Group(func(r chi.Router) {
			r.Use(auth, moderatorOnly)
			r.Get("/moderation/pending", moderationHandler.ListPending)
			r.Post("/moderation/{id}/approve", moderationHandler.Approve)
			r.Post("/moderation/{id}/reject", moderationHandler.Reject)
		})
	})

	return r
}
