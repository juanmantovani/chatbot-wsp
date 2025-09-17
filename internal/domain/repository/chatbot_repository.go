package repository

import (
	"sync"
	"time"

	"chatbot-wsp/internal/domain/errors"
	"chatbot-wsp/internal/domain/models"
)

// ChatbotRepository defines the interface for chatbot data operations
type ChatbotRepository interface {
	GetUserState(userID string) (*models.ChatbotState, error)
	SaveUserState(state *models.ChatbotState) error
	GetFlowByState(state string) (*models.ChatbotFlow, error)
	GetAllFlows() (map[string]*models.ChatbotFlow, error)
	StartSessionCleanup(expirationHours, cleanupIntervalMin int)
	StopSessionCleanup()
}

// InMemoryChatbotRepository implements ChatbotRepository using in-memory storage
type InMemoryChatbotRepository struct {
	userStates      map[string]*models.ChatbotState
	flows           map[string]*models.ChatbotFlow
	mutex           sync.RWMutex
	stopCleanup     chan bool
	expirationHours int
}

// NewInMemoryChatbotRepository creates a new in-memory repository
func NewInMemoryChatbotRepository() *InMemoryChatbotRepository {
	repo := &InMemoryChatbotRepository{
		userStates:      make(map[string]*models.ChatbotState),
		flows:           make(map[string]*models.ChatbotFlow),
		stopCleanup:     make(chan bool),
		expirationHours: 24, // Default to 24 hours
	}

	// Initialize default flows
	repo.initializeFlows()

	return repo
}

// GetUserState retrieves the current state of a user
func (r *InMemoryChatbotRepository) GetUserState(userID string) (*models.ChatbotState, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	state, exists := r.userStates[userID]
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
	if r.isSessionExpired(state) {
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

// SaveUserState saves the current state of a user
func (r *InMemoryChatbotRepository) SaveUserState(state *models.ChatbotState) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.userStates[state.UserID] = state
	return nil
}

// GetFlowByState retrieves the flow configuration for a given state
func (r *InMemoryChatbotRepository) GetFlowByState(state string) (*models.ChatbotFlow, error) {
	flow, exists := r.flows[state]
	if !exists {
		return nil, errors.ErrFlowNotFound
	}
	return flow, nil
}

// GetAllFlows retrieves all available flows
func (r *InMemoryChatbotRepository) GetAllFlows() (map[string]*models.ChatbotFlow, error) {
	return r.flows, nil
}

// initializeFlows sets up the default conversation flows
func (r *InMemoryChatbotRepository) initializeFlows() {
	// Welcome flow
	r.flows["welcome"] = &models.ChatbotFlow{
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

	// Option A flow - Consulta médica telefónica
	r.flows["option_a"] = &models.ChatbotFlow{
		State: "option_a",
		Message: `La consulta telefónica es un acto médico y tiene un valor de $15.000 ARS (no cubierta por obra social).
Para avanzar, enviá:
1️⃣ Nombre y edad del paciente
2️⃣ Motivo de la consulta
3️⃣ Comprobante de pago (Alias: Narvaez.Carla.B)

Información importante:
https://appar.com.ar/consulta-pediatrica-online/

📌 Una vez completados estos pasos, la Dra. se pondrá en contacto.`,
		DataRequest: "datos_consulta_medica",
	}

	// Option B flow - Lectura de estudios
	r.flows["option_b"] = &models.ChatbotFlow{
		State: "option_b",
		Message: `Por favor enviá:
1️⃣ Fotos claras o PDF de los estudios
2️⃣ Síntomas actuales y fecha de realización
3️⃣ Tu duda o pregunta principal
4️⃣ Comprobante de pago (Alias: Narvaez.Carla.B) $15.000 ARS

Información importante:
https://appar.com.ar/consulta-pediatrica-online/

📌 Una vez completados estos pasos, la Dra. se pondrá en contacto.`,
		DataRequest: "datos_lectura_estudios",
	}

	// Option C flow - Solicitar turno en consultorio
	r.flows["option_c"] = &models.ChatbotFlow{
		State: "option_c",
		Message: `Para turnos comunicarse a los siguientes números
– Centro Médico Cervantes (WhatsApp: 343-4066281)
– Consultorios OSPEP (WhatsApp: 343-5138637)`,
		DataRequest: "datos_turno",
	}

	// Option D flow - Información sobre BabyHome
	r.flows["option_d"] = &models.ChatbotFlow{
		State: "option_d",
		Message: `💜 ¡Qué alegría que te interese BabyHome!
Ofrecemos:
✅ Consulta prenatal
✅ Recepción neonatal personalizada (COPAP y primera hora siempre que mamá y bebé estén clínicamente bien)
✅ Controles en domicilio

Para orientarte, contanos:
1️⃣ Semana de embarazo / FPP
2️⃣ Maternidad y obstetra
3️⃣ Si desean priorizar COPAP/primera hora
4️⃣ Si quieren coordinar una consulta prenatal`,
		DataRequest: "datos_babyhome",
	}

	// Data collection flows
	r.flows["collecting_data"] = &models.ChatbotFlow{
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

	// Invalid option validation flow
	r.flows["invalid_option"] = &models.ChatbotFlow{
		State:   "invalid_option",
		Message: `⚠️ Por favor, ingresa una opción válida (A, B, C o D).`,
	}
}

// isSessionExpired checks if a session has expired based on UpdatedAt timestamp
func (r *InMemoryChatbotRepository) isSessionExpired(state *models.ChatbotState) bool {
	expirationDuration := time.Duration(r.expirationHours) * time.Hour
	return time.Since(state.UpdatedAt) > expirationDuration
}

// StartSessionCleanup starts the background cleanup goroutine
func (r *InMemoryChatbotRepository) StartSessionCleanup(expirationHours, cleanupIntervalMin int) {
	r.expirationHours = expirationHours // Update the expiration hours

	go func() {
		ticker := time.NewTicker(time.Duration(cleanupIntervalMin) * time.Minute)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				r.cleanupExpiredSessions()
			case <-r.stopCleanup:
				return
			}
		}
	}()
}

// StopSessionCleanup stops the background cleanup goroutine
func (r *InMemoryChatbotRepository) StopSessionCleanup() {
	close(r.stopCleanup)
}

// cleanupExpiredSessions removes expired sessions from memory
func (r *InMemoryChatbotRepository) cleanupExpiredSessions() {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	expirationDuration := time.Duration(r.expirationHours) * time.Hour
	now := time.Now()

	var expiredUsers []string
	for userID, state := range r.userStates {
		if now.Sub(state.UpdatedAt) > expirationDuration {
			expiredUsers = append(expiredUsers, userID)
		}
	}

	// Remove expired sessions
	for _, userID := range expiredUsers {
		delete(r.userStates, userID)
	}

	if len(expiredUsers) > 0 {
		// Log cleanup activity (you can add logging here if needed)
		// fmt.Printf("Cleaned up %d expired sessions\n", len(expiredUsers))
	}
}
