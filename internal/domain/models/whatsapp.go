package models

import "time"

// WhatsAppMessage represents a WhatsApp message
type WhatsAppMessage struct {
	ID        string    `json:"id"`
	From      string    `json:"from"`
	To        string    `json:"to"`
	Text      string    `json:"text"`
	Type      string    `json:"type"`
	Timestamp time.Time `json:"timestamp"`
}

// WhatsAppWebhook represents the webhook payload from WhatsApp
type WhatsAppWebhook struct {
	Object string `json:"object"`
	Entry  []struct {
		ID      string `json:"id"`
		Changes []struct {
			Value struct {
				MessagingProduct string `json:"messaging_product"`
				Metadata         struct {
					DisplayPhoneNumber string `json:"display_phone_number"`
					PhoneNumberID      string `json:"phone_number_id"`
				} `json:"metadata"`
				Messages []struct {
					From      string `json:"from"`
					ID        string `json:"id"`
					Timestamp string `json:"timestamp"`
					Text      struct {
						Body string `json:"body"`
					} `json:"text"`
					Type string `json:"type"`
				} `json:"messages"`
			} `json:"value"`
			Field string `json:"field"`
		} `json:"changes"`
	} `json:"entry"`
}

// WhatsAppResponse represents the response to send to WhatsApp
type WhatsAppResponse struct {
	MessagingProduct string `json:"messaging_product"`
	To               string `json:"to"`
	Type             string `json:"type"`
	Text             struct {
		Body string `json:"body"`
	} `json:"text"`
}

// ChatbotState represents the current state of a user conversation
type ChatbotState struct {
	UserID    string            `json:"user_id"`
	State     string            `json:"state"`
	Option    string            `json:"option,omitempty"`
	Data      map[string]string `json:"data,omitempty"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
}

// ChatbotOption represents a menu option
type ChatbotOption struct {
	ID          string `json:"id"`
	Label       string `json:"label"`
	Description string `json:"description"`
	NextState   string `json:"next_state"`
}

// ChatbotFlow represents the conversation flow
type ChatbotFlow struct {
	State       string          `json:"state"`
	Message     string          `json:"message"`
	Options     []ChatbotOption `json:"options,omitempty"`
	DataRequest string          `json:"data_request,omitempty"`
}
