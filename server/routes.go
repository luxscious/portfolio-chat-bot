package main

import (
	"encoding/json"
	"net/http"
	"time"

	"go-ai/config"
	"go-ai/db"
	"go-ai/openai"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

// ─────────────────────────────────────────────────────────────────────────────
// Request/Response Types
// ─────────────────────────────────────────────────────────────────────────────

type ChatRequest struct {
	UserID  string `json:"userId"`
	Message string `json:"content"`
}

type ChatResponse struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ─────────────────────────────────────────────────────────────────────────────
// POST /chat — handles user input and returns GPT response
// ─────────────────────────────────────────────────────────────────────────────

func chatHandler(w http.ResponseWriter, r *http.Request) {
	var req ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	reply, err := openai.SmartQuery(req.UserID, req.Message)
	if err != nil {
		http.Error(w, "Failed to generate response: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := storeChatPair(req.UserID, req.Message, reply); err != nil {
		http.Error(w, "Failed to store chat messages: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ChatResponse{
		Role:    "assistant",
		Content: reply,
	})
}

// ─────────────────────────────────────────────────────────────────────────────
// GET /chat?userId=... — fetches past chat messages
// ─────────────────────────────────────────────────────────────────────────────

func handleGetChat(w http.ResponseWriter, r *http.Request) {
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

// ─────────────────────────────────────────────────────────────────────────────
// Internal: Store both user + assistant message to DB
// ─────────────────────────────────────────────────────────────────────────────

func storeChatPair(userId, userMsg, assistantMsg string) error {
	now := time.Now()
	for _, msg := range []db.ChatMessage{
		{UserID: userId, Role: "user", Content: userMsg, Timestamp: now},
		{UserID: userId, Role: "assistant", Content: assistantMsg, Timestamp: now},
	} {
		if err := db.StoreMessage(userId, msg); err != nil {
			return err
		}
	}
	return nil
}

// ─────────────────────────────────────────────────────────────────────────────
// RegisterRoutes sets up HTTP routes and middleware
// ─────────────────────────────────────────────────────────────────────────────

func RegisterRoutes() http.Handler {
	r := chi.NewRouter()

	// Middleware stack
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{config.GetFrontendOrigin()},
		AllowedMethods:   []string{"GET", "POST"},
		AllowedHeaders:   []string{"Accept", "Content-Type"},
		AllowCredentials: true,
	}))

	// Routes
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("👋 Resume Chatbot backend is running!"))
	})
	r.Get("/chat", handleGetChat)
	r.Post("/chat", chatHandler)

	return r
}

