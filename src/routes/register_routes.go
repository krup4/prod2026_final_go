package routes

import (
	"backend/handlers"

	"github.com/gin-gonic/gin"
)

func registerUserRoutes(group *gin.RouterGroup, userHandler *handlers.UserHandler) {
	users := group.Group("/users")
	users.POST("", userHandler.Create)
	users.GET("/list", userHandler.List)
}

func registerExperimentsRoutes(group *gin.RouterGroup, experimentsHandler *handlers.ExperimentsHandler) {
	experiments := group.Group("/experiments")
	experiments.POST("", experimentsHandler.Create)
	experiments.PATCH("", experimentsHandler.ChangeStatus)
}

func registerFeatureRoutes(group *gin.RouterGroup, featureHandler *handlers.FeatureHandler) {
	feature := group.Group("/feature")
	feature.POST("", featureHandler.GetFeature)
}
