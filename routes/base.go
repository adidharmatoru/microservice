package routes

import (
	"microservice/routes/api"
	"microservice/routes/web"

	"github.com/gin-gonic/gin"
)

// SetupRouter sets up the routes for the Gin engine
func SetupRouter() *gin.Engine {
	router := gin.Default()

	// Setup base routes
	web.SetupBaseRoutes(router)

	// Setup V1 API routes
	api.SetupV1Routes(router)

	return router
}
