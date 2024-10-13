package api

import (
    "github.com/gin-gonic/gin"
    "microservice/routes/api/v1"
)

func SetupV1Routes(router *gin.Engine) {
    v1.SetupUserRoutes(router)
    v1.SetupAuthRoutes(router)
}
