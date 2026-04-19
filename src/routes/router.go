package routes

import (
	"backend/handlers"

	"github.com/gin-gonic/gin"
)

type Handlers struct {
	User        *handlers.UserHandler
	Experiments *handlers.ExperimentsHandler
	Feature     *handlers.FeatureHandler
}

func RegisterV1Routes(group *gin.RouterGroup, h Handlers) {
	registerUserRoutes(group, h.User)
	registerExperimentsRoutes(group, h.Experiments)
	registerFeatureRoutes(group, h.Feature)
}
