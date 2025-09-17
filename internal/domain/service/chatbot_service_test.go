package service

import (
	"strings"
	"testing"
	"time"

	"chatbot-wsp/internal/domain/errors"
	"chatbot-wsp/internal/domain/models"
)

// mockRepository is a mock implementation of ChatbotRepository
type mockRepository struct {
	userStates      map[string]*models.ChatbotState
	flows           map[string]*models.ChatbotFlow
	expirationHours int
}

func newMockRepository() *mockRepository {
	repo := &mockRepository{
		userStates:      make(map[string]*models.ChatbotState),
		flows:           make(map[string]*models.ChatbotFlow),
		expirationHours: 24, // Default to 24 hours
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
		Message:     "La consulta telefónica es un acto médico y tiene un valor de $15.000 ARS (no cubierta por obra social).",
		DataRequest: "datos_consulta_medica",
	}

	repo.flows["option_b"] = &models.ChatbotFlow{
		State:       "option_b",
		Message:     "Por favor enviá: Fotos claras o PDF de los estudios",
		DataRequest: "datos_lectura_estudios",
	}

	repo.flows["option_c"] = &models.ChatbotFlow{
		State:       "option_c",
		Message:     "Para turnos comunicarse a los siguientes números",
		DataRequest: "datos_turno",
	}

	repo.flows["option_d"] = &models.ChatbotFlow{
		State:       "option_d",
		Message:     "¡Qué alegría que te interese BabyHome!",
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

	repo.flows["invalid_option"] = &models.ChatbotFlow{
		State:   "invalid_option",
		Message: `⚠️ Por favor, ingresa una opción válida (A, B, C o D).`,
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

	// Check if session has expired (lazy cleanup)
	if m.isSessionExpired(state) {
		// Return a fresh state instead of the expired one
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

func (m *mockRepository) StartSessionCleanup(expirationHours, cleanupIntervalMin int) {
	m.expirationHours = expirationHours // Update the expiration hours for tests
}

func (m *mockRepository) StopSessionCleanup() {
	// Mock implementation - do nothing for tests
}

func (m *mockRepository) isSessionExpired(state *models.ChatbotState) bool {
	expirationDuration := time.Duration(m.expirationHours) * time.Hour
	return time.Since(state.UpdatedAt) > expirationDuration
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

func TestChatbotService_InvalidOptionValidation(t *testing.T) {
	repo := newMockRepository()
	service := NewChatbotService(repo)

	tests := []struct {
		name             string
		userID           string
		message          string
		expectedContains string
		expectedState    string
	}{
		{
			name:             "Invalid option in welcome state",
			userID:           "user123",
			message:          "X",
			expectedContains: "⚠️ Por favor, ingresa una opción válida (A, B, C o D).",
			expectedState:    "welcome",
		},
		{
			name:             "Invalid option with lowercase",
			userID:           "user123",
			message:          "x",
			expectedContains: "⚠️ Por favor, ingresa una opción válida (A, B, C o D).",
			expectedState:    "welcome",
		},
		{
			name:             "Invalid option with number",
			userID:           "user123",
			message:          "1",
			expectedContains: "⚠️ Por favor, ingresa una opción válida (A, B, C o D).",
			expectedState:    "welcome",
		},
		{
			name:             "Empty message",
			userID:           "user123",
			message:          "",
			expectedContains: "⚠️ Por favor, ingresa una opción válida (A, B, C o D).",
			expectedState:    "welcome",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, err := service.ProcessMessage(tt.userID, tt.message)

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if response == nil {
				t.Errorf("Expected response but got nil")
				return
			}

			// Verify the response contains the validation message
			if !strings.Contains(response.Text.Body, tt.expectedContains) {
				t.Errorf("Expected response to contain '%s', but got: %s", tt.expectedContains, response.Text.Body)
			}

			// Verify the response also contains the welcome message
			if !strings.Contains(response.Text.Body, "Chatbot BabyHome") {
				t.Errorf("Expected response to contain welcome message, but got: %s", response.Text.Body)
			}

			// Verify user state was updated correctly
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

func TestChatbotService_SessionExpiration(t *testing.T) {
	repo := newMockRepository()
	service := NewChatbotService(repo)

	// Create a user state with an old timestamp (simulating expired session)
	oldTime := time.Now().Add(-25 * time.Hour) // 25 hours ago
	userState := &models.ChatbotState{
		UserID:    "user123",
		State:     "collecting_data",
		Option:    "A",
		Data:      map[string]string{"test": "data"},
		CreatedAt: oldTime,
		UpdatedAt: oldTime,
	}
	repo.SaveUserState(userState)

	// Process a message - should reset to welcome state due to expiration
	response, err := service.ProcessMessage("user123", "hello")

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if response == nil {
		t.Errorf("Expected response but got nil")
		return
	}

	// Verify that the response contains the welcome message (indicating reset)
	if !strings.Contains(response.Text.Body, "Chatbot BabyHome") {
		t.Errorf("Expected response to contain welcome message, indicating session reset")
	}

	// Verify that the user state was reset
	newUserState, err := repo.GetUserState("user123")
	if err != nil {
		t.Errorf("Failed to get user state: %v", err)
	}

	if newUserState.State != "welcome" {
		t.Errorf("Expected state to be reset to 'welcome', got %s", newUserState.State)
	}

	if len(newUserState.Data) != 0 {
		t.Errorf("Expected data to be cleared, got %v", newUserState.Data)
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
