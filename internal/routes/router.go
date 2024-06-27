package routes

import (
	"chat-system/internal/auth"
	"chat-system/internal/message"
	"chat-system/internal/middleware"
	"github.com/gorilla/mux"
)

func SetupRouter() *mux.Router {
	router := mux.NewRouter()

	// Authentication routes
	router.HandleFunc("/register", auth.RegisterHandler).Methods("POST")
	router.HandleFunc("/login", auth.LoginHandler).Methods("POST")

	// Messaging routes with authentication middleware
	messageRouter := router.PathPrefix("/").Subrouter()
	messageRouter.Use(middleware.AuthMiddleware)
	messageRouter.HandleFunc("/send", message.SendMessageHandler).Methods("POST")
	messageRouter.HandleFunc("/messages", message.GetMessagesHandler).Methods("GET")

	return router
}
