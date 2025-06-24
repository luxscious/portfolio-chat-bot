package main

import (
	"go-ai/config"
	"go-ai/db"
	"log"
	"net/http"
)

func init() {
	// Load environment variables
	config.LoadEnv()

	// Initialize databases
	db.InitMongo()
	db.InitNeo4j()
}

func main() {
	port := config.GetServerPort()
	log.Printf("✅ Server started on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, RegisterRoutes()))
}
