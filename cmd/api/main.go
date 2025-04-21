package main

import (
	"log"

	"github.com/EngenMe/go-clean-arch-auth/internal/Infrastructure/repository/respositories"
	"github.com/EngenMe/go-clean-arch-auth/internal/core/features/mediator"
	"github.com/EngenMe/go-clean-arch-auth/internal/core/http"
	"github.com/EngenMe/go-clean-arch-auth/internal/useCase"
	"github.com/EngenMe/go-clean-arch-auth/pkg/config"
	"github.com/EngenMe/go-clean-arch-auth/pkg/database"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configs: %v", err)
	}

	db, err := database.NewPostgresDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	userRepo := respositories.NewUserRepository(db)
	userUseCase := useCase.NewUserUseCase(userRepo)
	med := mediator.NewMediator()

	med.RegisterUserHandlers(userUseCase)

	router := gin.Default()

	http.RegisterRoutes(router, med)

	log.Printf("Server starting on port %s", cfg.ServerPort)
	if err := router.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
