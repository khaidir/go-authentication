package routes

import (
	"auth-services/internal/handler"
	"auth-services/internal/http/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(r *gin.Engine, h handler.UserHandler) {
	public := r.Group("/api")
	{
		public.POST("/register", h.Register)
		public.POST("/login", h.Login)
	}

	private := r.Group("/api")
	private.Use(middleware.JWTAuthMiddleware())
	{
		private.GET("/profile/:id", h.Profile)
		// Add other protected endpoints here
	}
}
