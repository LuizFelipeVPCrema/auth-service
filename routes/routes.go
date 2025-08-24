package routes

import (
	"auth-service/handlers"
	"auth-service/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(authHandler *handlers.AuthHandler, authMiddleware *middleware.AuthMiddleware) *gin.Engine {
	router := gin.Default()

	// Middleware global
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Rotas públicas
	public := router.Group("/api/v1")
	{
		// Health check
		public.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"status":  "ok",
				"message": "Auth service is running",
			})
		})

		// Rotas de autenticação
		auth := public.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.RefreshToken)
			auth.POST("/validate", authHandler.ValidateToken)
		}

		// Rotas de clientes
		clients := public.Group("/clients")
		{
			clients.POST("/", authHandler.CreateClient)
		}
	}

	// Rotas protegidas
	protected := router.Group("/api/v1")
	protected.Use(authMiddleware.Authenticate())
	{
		// Rotas de autenticação protegidas
		auth := protected.Group("/auth")
		{
			auth.GET("/profile", authHandler.GetProfile)
		}
	}

	return router
}
