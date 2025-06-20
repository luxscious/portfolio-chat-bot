package main

import (
	"encoding/json"
	"fmt"
	"log"
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
	// Parse request body
	var req ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Fetch prior messages from MongoDB
	history, err := db.GetMessages(req.UserID)
	if err != nil {
		http.Error(w, "Failed to retrieve message history: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Inject relevant content
	retrievedChunks, err  := openai.SearchRelevantChunks(req.Message)
	log.Println("Retrived Chunks",retrievedChunks)
	if err != nil {
		http.Error(w, "Embedding error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	// Build the full message array (system + retrieved + history + user msg)
	messages := []db.ChatMessage{
		{
			Role:    "system",
			Content: systemPrompt,
		},
		{
			Role:    "assistant",
			Content: fmt.Sprintf("Here are some relevant parts of my resume that might help answer your question:\n\n%s\n\n", retrievedChunks),

		},
	}

	// Append conversation history and the latest user message
	messages = append(messages, history...)
	messages = append(messages, db.ChatMessage{
		Role:    "user",
		Content: req.Message,
	})


	// Call OpenAI with messages
	reply, err := openai.CallOpenAI(messages)
	if err != nil {
		http.Error(w, "OpenAI request failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Store both messages in MongoDB
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

	// Respond to frontend
	resp := ChatResponse{
		Role:    "assistant",
		Content: reply,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
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