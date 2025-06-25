package db

import (
	"context"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

// ─────────────────────────────────────────────────────────────────────────────
// PUBLIC FUNCTIONS
// ─────────────────────────────────────────────────────────────────────────────

// GetAllSkillsSorted returns all Skill nodes ordered by name.
func GetAllSkillsSorted() ([]Skill, error) {
	query := `
		MATCH (s:Skill)
		RETURN s
		ORDER BY s.name
	`
	return querySkills(query, nil)
}

// SearchSkillsByTag returns skills associated with a specific project tag.
func SearchSkillsByTag(tag string) ([]Skill, error) {
	query := `
		MATCH (s:Skill)<-[:USES]-(p:Project)-[:HAS_TAG]->(t:Tag {name: $tag})
		RETURN DISTINCT s.name AS name
		ORDER BY name
	`
	ctx := context.Background()
	session := Neo4jDriver.NewSession(ctx, neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeRead,
	})
	defer session.Close(ctx)

	result, err := session.Run(ctx, query, map[string]any{"tag": tag})
	if err != nil {
		return nil, err
	}

	var skills []Skill
	for result.Next(ctx) {
		record := result.Record()
		skills = append(skills, Skill{
			Name: asString(record, "name"),
		})
	}

	if err = result.Err(); err != nil {
		return nil, err
	}

	return skills, nil
}

// SearchSkillsByName performs a case-insensitive fuzzy match on skill name.
func SearchSkillsByName(name string) ([]Skill, error) {
	query := `
		MATCH (s:Skill)
		WHERE toLower(s.name) CONTAINS toLower($name)
		RETURN s
	`
	params := map[string]interface{}{"name": name}
	return querySkills(query, params)
}

// ─────────────────────────────────────────────────────────────────────────────
// INTERNAL HELPER
// ─────────────────────────────────────────────────────────────────────────────

// querySkills executes a generic Cypher query and returns Skill nodes.
func querySkills(cypher string, params map[string]interface{}) ([]Skill, error) {
	ctx := context.Background()
	session := Neo4jDriver.NewSession(ctx, neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeRead,
	})
	defer session.Close(ctx)

	result, err := session.Run(ctx, cypher, params)
	if err != nil {
		return nil, err
	}

	var skills []Skill
	for result.Next(ctx) {
		record := result.Record()
		node, ok := record.Get("s")
		if !ok {
			continue
		}

		props := node.(neo4j.Node).Props
		skill := Skill{
			Name: toString(props["name"]),
		}
		skills = append(skills, skill)
	}

	if err = result.Err(); err != nil {
		return nil, err
	}

	return skills, nil
}
