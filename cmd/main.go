package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"chatbot-wsp/internal/domain/repository"
	"chatbot-wsp/internal/domain/service"
	"chatbot-wsp/internal/infrastructure/config"
	"chatbot-wsp/internal/infrastructure/http/handlers"
	"chatbot-wsp/internal/infrastructure/http/routes"
	"chatbot-wsp/internal/infrastructure/logger"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize logger
	logger.Init(cfg.Logging.Level)
	log := logger.GetLogger()

	log.WithFields(map[string]interface{}{
		"port": cfg.Server.Port,
		"host": cfg.Server.Host,
	}).Info("Starting WhatsApp Chatbot service")

	// Initialize repository
	chatbotRepo := repository.NewInMemoryChatbotRepository()

	// Start session cleanup
	chatbotRepo.StartSessionCleanup(cfg.Session.ExpirationHours, cfg.Session.CleanupIntervalMin)
	defer chatbotRepo.StopSessionCleanup()

	// Initialize service
	chatbotService := service.NewChatbotService(chatbotRepo)

	// Initialize handler
	whatsappHandler := handlers.NewWhatsAppHandler(chatbotService, &handlers.Config{
		VerifyToken:   cfg.WhatsApp.VerifyToken,
		AccessToken:   cfg.WhatsApp.AccessToken,
		PhoneNumberID: cfg.WhatsApp.PhoneNumberID,
		MyPhoneNumber: cfg.WhatsApp.MyPhoneNumber,
	})

	// Setup routes
	router := routes.SetupRoutes(whatsappHandler)

	// Create HTTP server
	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port),
		Handler: router,
	}

	// Start server in a goroutine
	go func() {
		log.Info("Server starting...")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.WithError(err).Fatal("Failed to start server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down server...")

	// Give outstanding requests 30 seconds to complete
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.WithError(err).Fatal("Server forced to shutdown")
	}

	log.Info("Server exited")
}
