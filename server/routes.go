package main

import (
	"encoding/json"
	"net/http"
	"time"

	// for reading env
	"go-ai/config"
	"go-ai/db"
	"go-ai/openai"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

// STRUCTS
type ChatRequest struct {
	UserID  string `json:"userId"`
	Message string `json:"content"`
}

type ChatResponse struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func chatHandler(w http.ResponseWriter, r *http.Request) {
	var req ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	

	// Use SmartQuery: 2-step GPT process (filter → generate)
	reply, err := openai.SmartQuery(globalResumeData, req.Message)
	if err != nil {
		http.Error(w, "Failed to generate response: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Store messages in MongoDB
	err = db.StoreMessage(req.UserID, db.ChatMessage{
		UserID:    req.UserID,
		Role:      "user",
		Content:   req.Message,
		Timestamp: time.Now(),
	})
	if err != nil {
		http.Error(w, "Failed to store user message: "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = db.StoreMessage(req.UserID, db.ChatMessage{
		UserID:    req.UserID,
		Role:      "assistant",
		Content:   reply,
		Timestamp: time.Now(),
	})
	if err != nil {
		http.Error(w, "Failed to store assistant message: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ChatResponse{
		Role:    "assistant",
		Content: reply,
	})
}

func handleGetChat(w http.ResponseWriter, r *http.Request){
	userId := r.URL.Query().Get("userId")
    if userId == "" {
        http.Error(w, "Missing userId", http.StatusBadRequest)
        return
    }

    messages, err := db.GetMessages(userId)
    if err != nil {
        http.Error(w, "Failed to fetch messages", http.StatusInternalServerError)
        return
    }
	if messages == nil {
        messages = []db.ChatMessage{} 
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(messages)

}

func RegisterRoutes() http.Handler {
	r := chi.NewRouter()
	
	// 🔐 Global middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	origin := config.GetFrontendOrigin()
	// 🌐 Enable CORS
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{origin}, // React server
		AllowedMethods:   []string{"GET", "POST"},
		AllowedHeaders:   []string{"Accept", "Content-Type"},
		AllowCredentials: true,
	}))

	// ✅ Routes
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("👋 Resume Chatbot backend is running!"))
	})
	r.Get("/chat", handleGetChat)

	r.Post("/chat", chatHandler)

	return r
}