package main

import (
	"log"

	"backend/database"
	"backend/handlers"
	"backend/repository"
	"backend/routes"
	"backend/services"

	"github.com/gin-gonic/gin"
)

func main() {
	db, err := database.Connect()
	if err != nil {
		log.Fatalf("database connect error: %v", err)
	}

	if err := database.Migrate(db); err != nil {
		log.Fatalf("database migrate error: %v", err)
	}

	userRepo := repository.NewUserRepository(db)
	userService := services.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)

	healthService := services.NewHealthService()
	healthHandler := handlers.NewHealthHandler(healthService)

	experimentsRepo := repository.NewExperimentsRepository(db)
	experimentsService := services.NewExperimentsService(experimentsRepo)
	experimentsHandler := handlers.NewExperimentsHandler(experimentsService)

	featureService := services.NewFeatureService(userRepo, experimentsRepo)
	featureHandler := handlers.NewFeatureHandler(featureService)

	router := gin.Default()
	api := router.Group("/api")
	v1 := api.Group("/v1")
	routes.RegisterV1Routes(v1, routes.Handlers{
		User:        userHandler,
		Experiments: experimentsHandler,
		Feature:     featureHandler,
	})
	v1.GET("/ready", healthHandler.Ready)
	v1.GET("/health", healthHandler.Health)

	healthService.IsReady = true
	router.Run(":80")
}
