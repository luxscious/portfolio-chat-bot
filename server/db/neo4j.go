package db

import (
	"go-ai/config"
	"log"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

var Neo4jDriver neo4j.DriverWithContext

func InitNeo4j() {
	driver, err := neo4j.NewDriverWithContext(
		config.GetNeo4jURI(),
		neo4j.BasicAuth(config.GetNeo4jUser(), config.GetNeo4jPass(), ""),
	)
	if err != nil {
		log.Fatalf("Failed to connect to Neo4j: %v", err)
	}
	Neo4jDriver = driver
	log.Println("✅ Connected to Neo4j")
}
