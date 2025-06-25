package db

import (
	"context"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

// ─────────────────────────────────────────────────────────────────────────────
// PUBLIC FUNCTIONS
// ─────────────────────────────────────────────────────────────────────────────

// GetAllTagsSorted returns all tags sorted alphabetically.
func GetAllTagsSorted() ([]Tag, error) {
	query := `
		MATCH (t:Tag)
		RETURN t.name AS name
		ORDER BY name
	`
	return queryTags(query, nil)
}

// FindTagsBySkill returns tags linked to projects that use the given skill.
func FindTagsBySkill(skill string) ([]Tag, error) {
	query := `
		MATCH (t:Tag)<-[:HAS_TAG]-(p:Project)-[:USES]->(s:Skill {name: $skill})
		RETURN DISTINCT t.name AS name
		ORDER BY name
	`
	params := map[string]any{"skill": skill}
	return queryTags(query, params)
}

// ─────────────────────────────────────────────────────────────────────────────
// INTERNAL HELPER
// ─────────────────────────────────────────────────────────────────────────────

// queryTags executes a Cypher query and returns Tag nodes.
func queryTags(cypher string, params map[string]any) ([]Tag, error) {
	result, err := withReadSession(func(tx neo4j.ManagedTransaction) (any, error) {
		res, err := tx.Run(context.Background(), cypher, params)
		if err != nil {
			return nil, err
		}

		var tags []Tag
		for res.Next(context.Background()) {
			record := res.Record()
			tags = append(tags, Tag{
				Name: asString(record, "name"),
			})
		}
		return tags, res.Err()
	})

	if err != nil {
		return nil, err
	}

	return result.([]Tag), nil
}
