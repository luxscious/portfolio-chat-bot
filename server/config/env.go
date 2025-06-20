package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)
func LoadEnv() {
	err := godotenv.Load(".env") 
	if err != nil {
		log.Println("⚠️  No .env file found — relying on external environment variables.")
	} 
}
// GetOpenAIKey fetches the OpenAI API key from environment
func GetOpenAIKey() string {
	key := os.Getenv("OPENAI_API_KEY")
	if key == "" {
		log.Fatal("❌ OPENAI_API_KEY not set in environment")
	}
	return key
}

// GetOpenAIChatURL fetches the OpenAI Chat URL
func GetOpenAIChatURL() string {
	url := os.Getenv("OPENAI_API_URL")
	if url == "" {
		log.Fatal("❌ OPENAI_API_URL not set in environment")
	}
	return url
}

// GetOpenAIEmbeddingURL fetches the OpenAI Embedding URL
func GetOpenAIEmbeddingURL() string {
	url := os.Getenv("OPENAI_EMBEDDING_URL")
	if url == "" {
		log.Fatal("❌ OPENAI_EMBEDDING_URL not set in environment")
	}
	return url
}

// GetFrontendOrigin returns the allowed CORS origin
func GetFrontendOrigin() string {
	origin := os.Getenv("FRONTEND_ORIGIN")
	if origin == "" {
		log.Fatal("❌ FRONTEND_ORIGIN not set in environment")
	}
	return origin
}

// GetServerPort returns the port the server should run on
func GetServerPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // fallback to default
	}
	return ":" + port
}


// Mongo config accessors
func GetMongoURI() string {
	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		log.Fatal("❌ MONGO_URI not set in environment")
	}
	return uri
}

func GetMongoDB() string {
	db := os.Getenv("MONGO_DB")
	if db == "" {
		log.Fatal("❌ MONGO_DB not set in environment")
	}
	return db
}

func GetMongoCollection() string {
	coll := os.Getenv("MONGO_COLLECTION")
	if coll == "" {
		log.Fatal("❌ MONGO_COLLECTION not set in environment")
	}
	return coll
}