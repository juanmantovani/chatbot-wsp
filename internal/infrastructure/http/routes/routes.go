package routes

import (
	"chatbot-wsp/internal/infrastructure/http/handlers"
	"chatbot-wsp/internal/infrastructure/http/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all routes for the application
func SetupRoutes(whatsappHandler *handlers.WhatsAppHandler) *gin.Engine {
	// Set Gin to release mode for production
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	// Add middleware
	router.Use(middleware.LoggingMiddleware())
	router.Use(middleware.RecoveryMiddleware())
	router.Use(middleware.CORS())

	// Health check endpoint
	router.GET("/health", whatsappHandler.HealthCheck)

	// Stats endpoint
	router.GET("/stats", whatsappHandler.GetStats)

	// WhatsApp webhook endpoints
	whatsapp := router.Group("/whatsapp")
	{
		whatsapp.GET("/webhook", whatsappHandler.VerifyWebhook)
		whatsapp.POST("/webhook", whatsappHandler.HandleWebhook)
		whatsapp.GET("/welcome", whatsappHandler.GetWelcomeMessage)
	}

	// API v1 endpoints
	api := router.Group("/api/v1")
	{
		api.GET("/health", whatsappHandler.HealthCheck)
		api.GET("/stats", whatsappHandler.GetStats)
	}

	return router
}
