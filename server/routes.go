package main

import (
	"encoding/json"
	"net/http"
	"os" // for reading env

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

// STRUCTS
type ChatRequest struct {
	UserID  string `json:"userId"`
	Message string `json:"message"`
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

	// Retrieve previous conversation history
	history := memory.Get(req.UserID)

	// Build full message array with system prompt
	messages := []ChatMessage{
		{
			Role:    "system",
			Content: systemPrompt, // global variable
		},
	}
	messages = append(messages, history...)
	messages = append(messages, ChatMessage{
		Role:    "user",
		Content: req.Message,
	})

	// Call OpenAI
	assistantReply, err := callOpenAI(messages)
	if err != nil {
		http.Error(w, "Failed to call OpenAI: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Update conversation memory
	memory.Add(req.UserID, ChatMessage{Role: "user", Content: req.Message})
	memory.Add(req.UserID, ChatMessage{Role: "assistant", Content: assistantReply})

	// Respond with structured assistant message
	resp := ChatResponse{
		Role:    "assistant",
		Content: assistantReply,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}


func RegisterRoutes() http.Handler {
	r := chi.NewRouter()
	
	// 🔐 Global middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	origin := os.Getenv("FRONTEND_ORIGIN")
	if origin == "" {
		origin = "http://localhost:5173" // fallback
	}
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

	r.Post("/chat", chatHandler)

	return r
}