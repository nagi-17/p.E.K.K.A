package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/nagi-17/p.E.K.K.A/internal/controllers"
	"github.com/nagi-17/p.E.K.K.A/internal/middleware"
)

func InitRouter() *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.EnableCORS)

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
		protectedRoutes.Post("/village/upgrade/cancel", controllers.CancelUpgradeHandler)
		protectedRoutes.Post("/village/collect", controllers.CollectResourceHandler)

		protectedRoutes.Post("/village/lab/troop/upgrade/start", controllers.StartTroopUpgrade)
		protectedRoutes.Post("/village/lab/troop/upgrade/finish", controllers.FinishTroopUpgrade)
		protectedRoutes.Get("/village/army", controllers.GetArmy)
		protectedRoutes.Post("/village/army/train", controllers.TrainArmy)
		protectedRoutes.Post("/village/army/discard", controllers.DiscardArmy)

		protectedRoutes.Get("/battle/matchmake", controllers.MatchMakeHandler)
		protectedRoutes.Post("/battle/attack", controllers.AttackHandler)
	})

	return router
}
