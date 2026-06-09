package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/nagi-17/p.E.K.K.A/internal/controllers"
	"github.com/nagi-17/p.E.K.K.A/internal/middleware"
)

func InitRouter() *chi.Mux {
	router := chi.NewRouter()

	router.Post("/register", controllers.Register)
	router.Post("/login", controllers.Login)

	router.Group(func(protectedRoutes chi.Router) {
		protectedRoutes.Use(middleware.VerifyJWT)
	})

	return router
}
