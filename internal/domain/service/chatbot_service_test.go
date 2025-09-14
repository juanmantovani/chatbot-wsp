package service

import (
	"testing"
	"time"

	"chatbot-wsp/internal/domain/errors"
	"chatbot-wsp/internal/domain/models"
)

// mockRepository is a mock implementation of ChatbotRepository
type mockRepository struct {
	userStates map[string]*models.ChatbotState
	flows      map[string]*models.ChatbotFlow
}

func newMockRepository() *mockRepository {
	repo := &mockRepository{
		userStates: make(map[string]*models.ChatbotState),
		flows:      make(map[string]*models.ChatbotFlow),
	}

	// Initialize test flows
	repo.flows["welcome"] = &models.ChatbotFlow{
		State: "welcome",
		Message: `🤖 Chatbot BabyHome – Dra. Carla Narváez
👋 ¡Hola! Gracias por comunicarte.
Por favor, seleccioná una opción escribiendo la letra correspondiente:
A. Realizar consulta médica telefónica
B. Enviar estudios para lectura
C. Solicitar turno en consultorio
D. Consulta sobre BabyHome
(Si es una urgencia, por favor acudí a una guardia)`,
		Options: []models.ChatbotOption{
			{ID: "A", Label: "A", Description: "Realizar consulta médica telefónica", NextState: "option_a"},
			{ID: "B", Label: "B", Description: "Enviar estudios para lectura", NextState: "option_b"},
			{ID: "C", Label: "C", Description: "Solicitar turno en consultorio", NextState: "option_c"},
			{ID: "D", Label: "D", Description: "Consulta sobre BabyHome", NextState: "option_d"},
		},
	}

	repo.flows["option_a"] = &models.ChatbotFlow{
		State:       "option_a",
		Message:     "Consulta médica telefónica - $15.000 ARS",
		DataRequest: "datos_consulta_medica",
	}

	repo.flows["option_b"] = &models.ChatbotFlow{
		State:       "option_b",
		Message:     "Lectura de estudios - $15.000 ARS",
		DataRequest: "datos_lectura_estudios",
	}

	repo.flows["option_c"] = &models.ChatbotFlow{
		State:       "option_c",
		Message:     "Solicitar turno en consultorio",
		DataRequest: "datos_turno",
	}

	repo.flows["option_d"] = &models.ChatbotFlow{
		State:       "option_d",
		Message:     "Información sobre BabyHome",
		DataRequest: "datos_babyhome",
	}

	repo.flows["collecting_data"] = &models.ChatbotFlow{
		State: "collecting_data",
		Message: `Gracias por la información. ¿Hay algo más en lo que pueda ayudarte?
Por favor, seleccioná una opción escribiendo la letra correspondiente:
A. Realizar consulta médica telefónica
B. Enviar estudios para lectura
C. Solicitar turno en consultorio
D. Consulta sobre BabyHome`,
		Options: []models.ChatbotOption{
			{ID: "A", Label: "A", Description: "Realizar consulta médica telefónica", NextState: "option_a"},
			{ID: "B", Label: "B", Description: "Enviar estudios para lectura", NextState: "option_b"},
			{ID: "C", Label: "C", Description: "Solicitar turno en consultorio", NextState: "option_c"},
			{ID: "D", Label: "D", Description: "Consulta sobre BabyHome", NextState: "option_d"},
		},
	}

	return repo
}

func (m *mockRepository) GetUserState(userID string) (*models.ChatbotState, error) {
	state, exists := m.userStates[userID]
	if !exists {
		return &models.ChatbotState{
			UserID:    userID,
			State:     "welcome",
			Data:      make(map[string]string),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}, nil
	}
	return state, nil
}

func (m *mockRepository) SaveUserState(state *models.ChatbotState) error {
	m.userStates[state.UserID] = state
	return nil
}

func (m *mockRepository) GetFlowByState(state string) (*models.ChatbotFlow, error) {
	flow, exists := m.flows[state]
	if !exists {
		return nil, errors.ErrFlowNotFound
	}
	return flow, nil
}

func (m *mockRepository) GetAllFlows() (map[string]*models.ChatbotFlow, error) {
	return m.flows, nil
}

func TestChatbotService_ProcessMessage_WelcomeState(t *testing.T) {
	repo := newMockRepository()
	service := NewChatbotService(repo)

	tests := []struct {
		name          string
		userID        string
		message       string
		expectedState string
		expectError   bool
	}{
		{
			name:          "Valid option A",
			userID:        "user123",
			message:       "A",
			expectedState: "collecting_data",
			expectError:   false,
		},
		{
			name:          "Valid option B",
			userID:        "user123",
			message:       "B",
			expectedState: "collecting_data",
			expectError:   false,
		},
		{
			name:          "Invalid option",
			userID:        "user123",
			message:       "X",
			expectedState: "welcome",
			expectError:   false,
		},
		{
			name:          "Empty message",
			userID:        "user123",
			message:       "",
			expectedState: "welcome",
			expectError:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, err := service.ProcessMessage(tt.userID, tt.message)

			if tt.expectError && err == nil {
				t.Errorf("Expected error but got none")
			}

			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if response == nil {
				t.Errorf("Expected response but got nil")
				return
			}

			// Verify user state was updated
			userState, err := repo.GetUserState(tt.userID)
			if err != nil {
				t.Errorf("Failed to get user state: %v", err)
			}

			if userState.State != tt.expectedState {
				t.Errorf("Expected state %s, got %s", tt.expectedState, userState.State)
			}
		})
	}
}

func TestChatbotService_ProcessMessage_DataCollectionState(t *testing.T) {
	repo := newMockRepository()
	service := NewChatbotService(repo)

	// Set up user in collecting_data state
	userState := &models.ChatbotState{
		UserID:    "user123",
		State:     "collecting_data",
		Option:    "A",
		Data:      make(map[string]string),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	repo.SaveUserState(userState)

	tests := []struct {
		name          string
		message       string
		expectedState string
		expectError   bool
	}{
		{
			name:          "Select another option A",
			message:       "A",
			expectedState: "collecting_data",
			expectError:   false,
		},
		{
			name:          "Select option B",
			message:       "B",
			expectedState: "collecting_data",
			expectError:   false,
		},
		{
			name:          "Invalid option",
			message:       "X",
			expectedState: "welcome",
			expectError:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, err := service.ProcessMessage("user123", tt.message)

			if tt.expectError && err == nil {
				t.Errorf("Expected error but got none")
			}

			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if response == nil {
				t.Errorf("Expected response but got nil")
				return
			}

			// Verify user state was updated
			userState, err := repo.GetUserState("user123")
			if err != nil {
				t.Errorf("Failed to get user state: %v", err)
			}

			if userState.State != tt.expectedState {
				t.Errorf("Expected state %s, got %s", tt.expectedState, userState.State)
			}
		})
	}
}

func TestChatbotService_GetWelcomeMessage(t *testing.T) {
	repo := newMockRepository()
	service := NewChatbotService(repo)

	response := service.GetWelcomeMessage()

	if response == nil {
		t.Errorf("Expected welcome message but got nil")
		return
	}

	if response.MessagingProduct != "whatsapp" {
		t.Errorf("Expected messaging product 'whatsapp', got '%s'", response.MessagingProduct)
	}

	if response.Type != "text" {
		t.Errorf("Expected type 'text', got '%s'", response.Type)
	}

	if response.Text.Body == "" {
		t.Errorf("Expected non-empty message body")
	}
}

func TestIsValidOption(t *testing.T) {
	tests := []struct {
		option   string
		expected bool
	}{
		{"A", true},
		{"B", true},
		{"C", true},
		{"D", true},
		{"a", false},
		{"X", false},
		{"1", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.option, func(t *testing.T) {
			result := isValidOption(tt.option)
			if result != tt.expected {
				t.Errorf("isValidOption(%s) = %v, expected %v", tt.option, result, tt.expected)
			}
		})
	}
}
