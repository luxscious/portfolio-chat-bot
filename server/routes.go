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
	var req ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Get full thread (system + past + user message)
	history := memory.Get(req.UserID)
	system := ChatMessage{
		Role:    "system",
		Content: systemPrompt,
	}
	thread := append([]ChatMessage{system}, history...)
	thread = append(thread, ChatMessage{Role: "user", Content: req.Message})

	// Call OpenAI
	reply, err := callOpenAI(thread)
	if err != nil {
		http.Error(w, "OpenAI error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Save both user message and assistant reply
	memory.Add(req.UserID, ChatMessage{Role: "user", Content: req.Message})
	memory.Add(req.UserID, ChatMessage{Role: "assistant", Content: reply})

	// Return assistant reply
	json.NewEncoder(w).Encode(ChatResponse{Role: "assistant", Content: reply})
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