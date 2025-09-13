package routes

import (
	"auth-service/handlers"
	"auth-service/middleware"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(authHandler *handlers.AuthHandler, authMiddleware *middleware.AuthMiddleware) *gin.Engine {
	router := gin.Default()

	// Configuração do CORS
	corsConfig := cors.Config{
		AllowOrigins:     []string{"http://localhost:4000", "http://127.0.0.1:4000", "http://localhost:4200", "http://127.0.0.1:4200"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization", "X-Requested-With", "X-Client-ID"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	router.Use(cors.New(corsConfig))

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

		public.POST("/register", authHandler.Register)
		public.POST("/login", authHandler.Login)
		public.POST("/refresh", authHandler.RefreshToken)
		public.POST("/validate", authHandler.ValidateToken)

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
