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

		protectedRoutes.Get("/user/load", controllers.LoadPlayerInfo)

		protectedRoutes.Get("/village", controllers.LoadVillage)
		protectedRoutes.Post("/village/build", controllers.PlaceBuilding)
		protectedRoutes.Put("/village/move", controllers.MoveBuildingHandler)

		protectedRoutes.Post("/village/upgrade/start", controllers.StartUpgradeHandler)
		protectedRoutes.Post("/village/upgrade/finish", controllers.FinishUpgradeHandler)
		protectedRoutes.Post("/village/collect", controllers.CollectResourceHandler)
	})

	return router
}
