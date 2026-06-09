package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/nagi-17/p.E.K.K.A/internal/controllers"
)

func InitRouter() *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	router.Post("/register", controllers.Register)
	router.Post("/login", controllers.Login)

	return router
}
