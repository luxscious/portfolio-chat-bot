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

	// Load graph schema once at startup
	if err := db.LoadGraphSchemaOnce(); err != nil {
		log.Fatalf("❌ Failed to load graph schema: %v", err)
	}
	log.Println("✅ Graph schema loaded")
}

func main() {
	port := config.GetServerPort()
	addr := "0.0.0.0" + port

        log.Printf("✅ Server started on http://%s\n", addr)
        log.Fatal(http.ListenAndServe(addr, RegisterRoutes()))
}
