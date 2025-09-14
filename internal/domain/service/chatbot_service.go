package service

import (
	"fmt"
	"strings"
	"time"

	"chatbot-wsp/internal/domain/models"
	"chatbot-wsp/internal/domain/repository"
)

// ChatbotService defines the interface for chatbot business logic
type ChatbotService interface {
	ProcessMessage(userID, message string) (*models.WhatsAppResponse, error)
	GetWelcomeMessage() *models.WhatsAppResponse
}

// chatbotService implements ChatbotService
type chatbotService struct {
	repo repository.ChatbotRepository
}

// NewChatbotService creates a new chatbot service
func NewChatbotService(repo repository.ChatbotRepository) ChatbotService {
	return &chatbotService{
		repo: repo,
	}
}

// ProcessMessage processes incoming messages and returns appropriate responses
func (s *chatbotService) ProcessMessage(userID, message string) (*models.WhatsAppResponse, error) {
	// Get current user state
	userState, err := s.repo.GetUserState(userID)
	if err != nil {
		return nil, err
	}

	// Process message based on current state
	response, newState, err := s.processMessageByState(userState, message)
	if err != nil {
		return nil, err
	}

	// Update user state
	userState.State = newState
	userState.UpdatedAt = time.Now()
	if err := s.repo.SaveUserState(userState); err != nil {
		return nil, err
	}

	return response, nil
}

// processMessageByState handles message processing based on current state
func (s *chatbotService) processMessageByState(userState *models.ChatbotState, message string) (*models.WhatsAppResponse, string, error) {
	message = strings.TrimSpace(strings.ToUpper(message))

	switch userState.State {
	case "welcome":
		return s.handleWelcomeState(userState, message)
	case "option_a", "option_b", "option_c", "option_d":
		return s.handleOptionState(userState, message)
	case "collecting_data":
		return s.handleDataCollectionState(userState, message)
	default:
		return s.handleWelcomeState(userState, message)
	}
}

// handleWelcomeState processes messages in welcome state
func (s *chatbotService) handleWelcomeState(userState *models.ChatbotState, message string) (*models.WhatsAppResponse, string, error) {
	// Check if message is a valid option
	if isValidOption(message) {
		userState.Option = message
		flow, err := s.repo.GetFlowByState("option_" + strings.ToLower(message))
		if err != nil {
			return nil, "", err
		}

		response := &models.WhatsAppResponse{
			MessagingProduct: "whatsapp",
			To:               userState.UserID,
			Type:             "text",
		}
		response.Text.Body = flow.Message

		return response, "collecting_data", nil
	}

	// If not a valid option, show welcome message again
	flow, err := s.repo.GetFlowByState("welcome")
	if err != nil {
		return nil, "", err
	}

	response := &models.WhatsAppResponse{
		MessagingProduct: "whatsapp",
		To:               userState.UserID,
		Type:             "text",
	}
	response.Text.Body = s.formatWelcomeMessage(flow)

	return response, "welcome", nil
}

// handleOptionState processes messages when user has selected an option
func (s *chatbotService) handleOptionState(userState *models.ChatbotState, message string) (*models.WhatsAppResponse, string, error) {
	// Collect data based on the selected option
	dataKey := s.getDataKeyForOption(userState.Option)
	userState.Data[dataKey] = message

	// Move to collecting data state
	flow, err := s.repo.GetFlowByState("collecting_data")
	if err != nil {
		return nil, "", err
	}

	response := &models.WhatsAppResponse{
		MessagingProduct: "whatsapp",
		To:               userState.UserID,
		Type:             "text",
	}
	response.Text.Body = s.formatDataCollectionMessage(flow, userState)

	return response, "collecting_data", nil
}

// handleDataCollectionState processes messages after data collection
func (s *chatbotService) handleDataCollectionState(userState *models.ChatbotState, message string) (*models.WhatsAppResponse, string, error) {
	// Check if user wants to select another option
	if isValidOption(message) {
		userState.Option = message
		flow, err := s.repo.GetFlowByState("option_" + strings.ToLower(message))
		if err != nil {
			return nil, "", err
		}

		response := &models.WhatsAppResponse{
			MessagingProduct: "whatsapp",
			To:               userState.UserID,
			Type:             "text",
		}
		response.Text.Body = flow.Message

		return response, "collecting_data", nil
	}

	// Show menu again
	flow, err := s.repo.GetFlowByState("welcome")
	if err != nil {
		return nil, "", err
	}

	response := &models.WhatsAppResponse{
		MessagingProduct: "whatsapp",
		To:               userState.UserID,
		Type:             "text",
	}
	response.Text.Body = s.formatWelcomeMessage(flow)

	return response, "welcome", nil
}

// GetWelcomeMessage returns the initial welcome message
func (s *chatbotService) GetWelcomeMessage() *models.WhatsAppResponse {
	flow, _ := s.repo.GetFlowByState("welcome")
	if flow == nil {
		return &models.WhatsAppResponse{
			MessagingProduct: "whatsapp",
			Type:             "text",
			Text: struct {
				Body string `json:"body"`
			}{
				Body: "¡Hola! Bienvenido a nuestro servicio.",
			},
		}
	}

	return &models.WhatsAppResponse{
		MessagingProduct: "whatsapp",
		Type:             "text",
		Text: struct {
			Body string `json:"body"`
		}{
			Body: s.formatWelcomeMessage(flow),
		},
	}
}

// Helper functions

func isValidOption(message string) bool {
	validOptions := []string{"A", "B", "C", "D"}
	for _, option := range validOptions {
		if message == option {
			return true
		}
	}
	return false
}

func (s *chatbotService) getDataKeyForOption(option string) string {
	switch option {
	case "A":
		return "datos_consulta_medica"
	case "B":
		return "datos_lectura_estudios"
	case "C":
		return "datos_turno"
	case "D":
		return "datos_babyhome"
	default:
		return "data"
	}
}

func (s *chatbotService) formatWelcomeMessage(flow *models.ChatbotFlow) string {
	// El mensaje ya está formateado correctamente en el repositorio
	return flow.Message
}

func (s *chatbotService) formatDataCollectionMessage(flow *models.ChatbotFlow, userState *models.ChatbotState) string {
	// El mensaje ya está formateado correctamente en el repositorio
	message := flow.Message

	// Add collected data summary
	if len(userState.Data) > 0 {
		message += "\n\n📋 Datos recopilados:\n"
		for key, value := range userState.Data {
			message += fmt.Sprintf("• %s: %s\n", key, value)
		}
	}

	return message
}
