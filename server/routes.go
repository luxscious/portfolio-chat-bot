package main

import (
	"net/http"
	"os" // for reading env

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)
func chatHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("🔮 This will eventually return a chatbot response"))
}

func NewRouter() http.Handler {
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