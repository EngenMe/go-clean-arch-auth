package http

import (
	"github.com/EngenMe/go-clean-arch-auth/internal/core/features/mediator"
	"github.com/EngenMe/go-clean-arch-auth/internal/core/http/handler"
	"github.com/EngenMe/go-clean-arch-auth/internal/core/http/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, mediator *mediator.Mediator) {
	userHandler := handler.NewUserHandler(mediator)

	router.POST("/register", userHandler.Register)
	router.POST("/login", userHandler.Login)

	protected := router.Group("/users")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("/:id", userHandler.GetUser)
		protected.PUT("/:id", userHandler.UpdateUser)
		protected.DELETE("/:id", userHandler.DeleteUser)
		protected.GET("", userHandler.ListUsers)
	}
}
