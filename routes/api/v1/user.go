package v1

import (
	v1 "microservice/controllers/api/v1"
	"microservice/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupUserRoutes(router *gin.Engine) {
	users := router.Group("/api/v1/users")

	// Public routes
	users.OPTIONS("", v1.OptionsUsers)
	users.HEAD("", v1.HeadUsers)
	users.GET("", v1.ListUsers)
	users.GET("/:id", v1.GetUser)
	users.GET("/dummy", v1.DummyListUsers)

	// Private routes
	users.Use(middlewares.JWTMiddleware())
	{
		users.POST("", middlewares.CheckScope("create:users"), v1.CreateUser)
		users.PUT("/:id", middlewares.CheckScope("update:users"), v1.UpdateUser)
		users.PATCH("/:id", middlewares.CheckScope("update:users"), v1.UpdateUser)
		users.DELETE("/:id", middlewares.CheckScope("delete:users"), v1.DeleteUser)
	}
}
