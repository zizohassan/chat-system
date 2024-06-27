package message

import (
	"chat-system/internal/middleware"
	"chat-system/internal/models"
	"encoding/json"
	"net/http"
)

// MessageRequest represents the request body for sending a message
type MessageRequest struct {
	Recipient string `json:"recipient"`
	Content   string `json:"content"`
}

type MessageResponse struct {
	Message models.Message `json:"message"`
}

// SendMessageHandler handles sending a message
// @Summary Send a message
// @Description Send a message from one user to another
// @Tags message
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param message body MessageRequest true "Message content"
// @Success 200 {object} models.Message
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /send [post]
func SendMessageHandler(w http.ResponseWriter, r *http.Request) {
	username, ok := middleware.GetUsernameFromContext(r.Context())
	if !ok {
		http.Error(w, "unable to retrieve user from context", http.StatusInternalServerError)
		return
	}

	var req MessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	message, err := SendMessage(username, req.Recipient, req.Content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := MessageResponse{
		Message: message,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetMessagesHandler handles retrieving message history
// @Summary Get message history
// @Description Get the message history for a user
// @Tags message
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param username query string true "Username"
// @Success 200 {array} models.Message
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /messages [get]
func GetMessagesHandler(w http.ResponseWriter, r *http.Request) {
	sername, ok := middleware.GetUsernameFromContext(r.Context())
	if !ok {
		http.Error(w, "unable to retrieve user from context", http.StatusInternalServerError)
		return
	}

	messages, err := GetMessages(sername)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}
