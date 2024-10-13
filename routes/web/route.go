package web

import (
    "github.com/gin-gonic/gin"
    "github.com/swaggo/files"
    "github.com/swaggo/gin-swagger"
)

func SetupBaseRoutes(router *gin.Engine) {
    // Swagger documentation route
    router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
