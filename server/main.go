package main

import (
	"go-ai/config"
	"go-ai/db"
	"go-ai/openai"
	"log"
	"net/http"
)

func init() {
	// Load environment variables
	config.LoadEnv()

	// Initialize databases
	db.InitMongo()
	db.InitNeo4j()

	// (Optional) System prompt setup from static persona context
	systemPrompt, err := openai.BuildSystemPrompt("You are a helpful resume chatbot. Answer based on the user’s past experience.")
	if err != nil {
		log.Fatalf("❌ Failed to build system prompt: %v", err)
	}
	log.Println("✅ System prompt initialized")
	log.Println(systemPrompt)
}

func main() {
	port := config.GetServerPort()
	log.Printf("✅ Server started on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, RegisterRoutes()))
}
