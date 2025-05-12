package routes

import (
	"github.com/gin-gonic/gin"

	"auth-services/config"
	"auth-services/internal/handler"
	"auth-services/internal/repository"
	"auth-services/internal/usecase"
)

func InitRoutes(r *gin.Engine) {
	// Inisialisasi repository
	userRepo := repository.NewUserRepository(config.DB)

	// Inisialisasi usecase
	userUc := usecase.NewUserUsecase(userRepo)

	// Inisialisasi handler
	userHandler := handler.UserHandler{Uc: userUc}

	// Daftarkan seluruh routes
	RegisterAuthRoutes(r, userHandler)
}
