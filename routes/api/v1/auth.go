package v1

import (
	v1 "microservice/controllers/api/v1"

	"github.com/gin-gonic/gin"
)

func SetupAuthRoutes(router *gin.Engine) {
	auth := router.Group("/api/v1/oauth")

	// Public routes
	auth.POST("/token", v1.GenerateJWT)
}
