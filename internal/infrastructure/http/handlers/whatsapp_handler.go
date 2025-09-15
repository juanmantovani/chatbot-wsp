package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"chatbot-wsp/internal/domain/models"
	"chatbot-wsp/internal/domain/service"
	"chatbot-wsp/internal/infrastructure/logger"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// WhatsAppHandler handles WhatsApp webhook requests
type WhatsAppHandler struct {
	chatbotService service.ChatbotService
	config         *Config
}

// Config holds configuration for the handler
type Config struct {
	VerifyToken   string
	AccessToken   string
	PhoneNumberID string
	MyPhoneNumber string
}

// NewWhatsAppHandler creates a new WhatsApp handler
func NewWhatsAppHandler(chatbotService service.ChatbotService, config *Config) *WhatsAppHandler {
	return &WhatsAppHandler{
		chatbotService: chatbotService,
		config:         config,
	}
}

// VerifyWebhook handles WhatsApp webhook verification
func (h *WhatsAppHandler) VerifyWebhook(c *gin.Context) {
	mode := c.Query("hub.mode")
	token := c.Query("hub.verify_token")
	challenge := c.Query("hub.challenge")

	logger.GetLogger().WithFields(logrus.Fields{
		"mode":   mode,
		"token":  token,
		"action": "webhook_verification",
	}).Info("Webhook verification request")

	// Check if mode and token are correct
	if mode != "subscribe" || token != h.config.VerifyToken {
		logger.GetLogger().WithFields(logrus.Fields{
			"expected_token": h.config.VerifyToken,
			"received_token": token,
			"mode":           mode,
		}).Error("Webhook verification failed")

		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}

	// Respond with the challenge
	logger.GetLogger().Info("Webhook verification successful")
	c.String(http.StatusOK, challenge)
}

// HandleWebhook handles incoming WhatsApp messages
func (h *WhatsAppHandler) HandleWebhook(c *gin.Context) {
	var webhook models.WhatsAppWebhook

	if err := c.ShouldBindJSON(&webhook); err != nil {
		logger.GetLogger().WithError(err).Error("Failed to parse webhook payload")
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"error":  "Invalid JSON",
		})
		return
	}

	logger.GetLogger().WithFields(logrus.Fields{
		"object":  webhook.Object,
		"entries": len(webhook.Entry),
	}).Info("Received webhook")

	// Track processing results
	var processedMessages int
	var errors []string
	var totalMessages int
	var chatbotResponses []string

	// Process each entry
	for _, entry := range webhook.Entry {
		for _, change := range entry.Changes {
			if change.Field == "messages" {
				// Convert webhook messages to our message format
				var messages []models.WhatsAppMessage
				for _, msg := range change.Value.Messages {
					messages = append(messages, models.WhatsAppMessage{
						ID:   msg.ID,
						From: msg.From,
						Text: msg.Text.Body,
						Type: msg.Type,
					})
					totalMessages++
				}

				// Process messages and collect results
				processed, processingErrors, responses := h.processMessages(messages)
				processedMessages += processed
				errors = append(errors, processingErrors...)
				chatbotResponses = append(chatbotResponses, responses...)
			}
		}
	}

	// Prepare response based on processing results
	response := gin.H{
		"status":             "success",
		"messages_received":  totalMessages,
		"messages_processed": processedMessages,
	}

	if len(errors) > 0 {
		response["errors"] = errors
		response["status"] = "partial_success"
		logger.GetLogger().WithField("errors", errors).Warn("Some messages failed to process")
	}

	// Include chatbot responses for testing purposes
	if len(chatbotResponses) > 0 {
		response["chatbot_responses"] = chatbotResponses
	}

	c.JSON(http.StatusOK, response)
}

// processMessages processes incoming messages and returns processing statistics
func (h *WhatsAppHandler) processMessages(messages []models.WhatsAppMessage) (processed int, errors []string, responses []string) {
	for _, message := range messages {
		logger.GetLogger().WithFields(logrus.Fields{
			"from":       message.From,
			"message_id": message.ID,
			"type":       message.Type,
		}).Info("Processing message")

		// Only process text messages
		if message.Type != "text" {
			logger.GetLogger().WithField("type", message.Type).Warn("Ignoring non-text message")
			errors = append(errors, fmt.Sprintf("Message %s: unsupported type %s", message.ID, message.Type))
			continue
		}

		// Process the message
		response, err := h.chatbotService.ProcessMessage(message.From, message.Text)
		if err != nil {
			errorMsg := fmt.Sprintf("Message %s: failed to process - %v", message.ID, err)
			logger.GetLogger().WithError(err).Error("Failed to process message")
			errors = append(errors, errorMsg)
			continue
		}

		// Add response to the list for testing purposes
		if response != nil && response.Text.Body != "" {
			responses = append(responses, response.Text.Body)
		}

		// Send response back to WhatsApp
		if err := h.sendMessage(response); err != nil {
			errorMsg := fmt.Sprintf("Message %s: failed to send response - %v", message.ID, err)
			logger.GetLogger().WithError(err).Error("Failed to send response")
			errors = append(errors, errorMsg)
			continue
		}

		processed++
		logger.GetLogger().WithFields(logrus.Fields{
			"message_id": message.ID,
			"from":       message.From,
		}).Info("Message processed successfully")
	}

	return processed, errors, responses
}

// sendMessage sends a message to WhatsApp Business API
func (h *WhatsAppHandler) sendMessage(response *models.WhatsAppResponse) error {
	// Check if we have the required configuration
	if h.config.AccessToken == "" || h.config.PhoneNumberID == "" || h.config.MyPhoneNumber == "" {
		logger.GetLogger().Warn("WhatsApp configuration missing - skipping message send")
		return fmt.Errorf("WhatsApp configuration incomplete")
	}

	// Prepare the request payload
	payload := map[string]interface{}{
		"messaging_product": response.MessagingProduct,
		"to":                h.config.MyPhoneNumber,
		"type":              response.Type,
		"text": map[string]interface{}{
			"body": response.Text.Body,
		},
	}

	// Convert to JSON
	jsonData, err := json.Marshal(payload)
	if err != nil {
		logger.GetLogger().WithError(err).Error("Failed to marshal WhatsApp message payload")
		return fmt.Errorf("failed to marshal payload: %v", err)
	}

	// Create HTTP request
	url := fmt.Sprintf("https://graph.facebook.com/v17.0/%s/messages", h.config.PhoneNumberID)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		logger.GetLogger().WithError(err).Error("Failed to create WhatsApp API request")
		return fmt.Errorf("failed to create request: %v", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+h.config.AccessToken)

	// Send request
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		logger.GetLogger().WithError(err).Error("Failed to send message to WhatsApp API")
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.GetLogger().WithError(err).Error("Failed to read WhatsApp API response")
		return fmt.Errorf("failed to read response: %v", err)
	}

	// Log the response
	logger.GetLogger().WithFields(logrus.Fields{
		"status_code": resp.StatusCode,
		"response":    string(body),
		"to":          response.To,
		"type":        response.Type,
		"body":        response.Text.Body,
	}).Info("WhatsApp API response")

	// Check if the request was successful
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		logger.GetLogger().Info("Message sent successfully to WhatsApp")
		return nil
	}

	// Log error response
	logger.GetLogger().WithFields(logrus.Fields{
		"status_code": resp.StatusCode,
		"response":    string(body),
	}).Error("WhatsApp API returned error")

	return fmt.Errorf("WhatsApp API error: status %d, response: %s", resp.StatusCode, string(body))
}

// GetWelcomeMessage returns the welcome message
func (h *WhatsAppHandler) GetWelcomeMessage(c *gin.Context) {
	response := h.chatbotService.GetWelcomeMessage()
	c.JSON(http.StatusOK, response)
}

// HealthCheck returns the health status of the service
func (h *WhatsAppHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"timestamp": time.Now().UTC(),
		"service":   "chatbot-wsp",
	})
}

// GetStats returns basic statistics about the service
func (h *WhatsAppHandler) GetStats(c *gin.Context) {
	// In a real implementation, you would collect actual statistics
	stats := gin.H{
		"uptime":             time.Since(time.Now()).String(),
		"messages_processed": 0, // This would be tracked in a real implementation
		"active_users":       0, // This would be tracked in a real implementation
		"timestamp":          time.Now().UTC(),
	}

	c.JSON(http.StatusOK, stats)
}
