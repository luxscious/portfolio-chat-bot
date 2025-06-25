package db

import (
	"context"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

// ─────────────────────────────────────────────────────────────────────────────
// PUBLIC FUNCTIONS
// ─────────────────────────────────────────────────────────────────────────────

// GetAllHobbies returns all Hobby nodes sorted by name.
func GetAllHobbies() ([]Hobby, error) {
	query := `
		MATCH (h:Hobby)
		RETURN h
		ORDER BY h.name
	`
	return queryHobbies(query, nil)
}

// SearchHobbiesByName returns hobbies where the name partially matches the input (case-insensitive).
func SearchHobbiesByName(name string) ([]Hobby, error) {
	query := `
		MATCH (h:Hobby)
		WHERE toLower(h.name) CONTAINS toLower($name)
		RETURN h
	`
	params := map[string]interface{}{"name": name}
	return queryHobbies(query, params)
}

// FindHobbiesByTag returns hobbies associated with a specific tag.
func SearchHobbiesByTag(tag string) ([]Hobby, error) {
	query := `
		MATCH (h:Hobby)-[:HAS_TAG]->(t:Tag {name: $tag})
		RETURN h
	`
	params := map[string]interface{}{"tag": tag}
	return queryHobbies(query, params)
}

// ─────────────────────────────────────────────────────────────────────────────
// INTERNAL HELPER
// ─────────────────────────────────────────────────────────────────────────────

// queryHobbies runs a Cypher query and returns a slice of Hobby nodes.
func queryHobbies(cypher string, params map[string]interface{}) ([]Hobby, error) {
	ctx := context.Background()
	session := Neo4jDriver.NewSession(ctx, neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeRead,
	})
	defer session.Close(ctx)

	result, err := session.Run(ctx, cypher, params)
	if err != nil {
		return nil, err
	}

	var hobbies []Hobby
	for result.Next(ctx) {
		record := result.Record()
		node, ok := record.Get("h")
		if !ok {
			continue
		}

		props := node.(neo4j.Node).Props
		hobby := Hobby{
			Name:        toString(props["name"]),
			Description: toString(props["description"]),
		}
		hobbies = append(hobbies, hobby)
	}

	if err := result.Err(); err != nil {
		return nil, err
	}

	return hobbies, nil
}
