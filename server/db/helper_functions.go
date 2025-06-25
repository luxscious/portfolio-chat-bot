package db

import (
	"context"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func parseProjectNode(p any) Project {
	props := p.(neo4j.Node).Props
	return Project{
		ID:            props["id"].(string),
		Name:          props["name"].(string),
		Description:   props["description"].(string),
		Institution:   props["institution"].(string),
		Image:         props["image"].(string),
		Featured:      props["featured"].(bool),
		Contributions: toStringSlice(props["contributions"]),
		StartDate:     props["startDate"].(string),
		EndDate:       props["endDate"].(string),
		Demo:          safeString(props["demo"]),
		GitHub:        safeString(props["github"]),
	}
}

func parseSkillList(val any) []Skill {
	nodes, ok := val.([]any)
	if !ok {
		return nil
	}
	var skills []Skill
	for _, n := range nodes {
		skills = append(skills, Skill{Name: n.(neo4j.Node).Props["name"].(string)})
	}
	return skills
}

func parseTagList(val any) []Tag {
	nodes, ok := val.([]any)
	if !ok {
		return nil
	}
	var tags []Tag
	for _, n := range nodes {
		tags = append(tags, Tag{Name: n.(neo4j.Node).Props["name"].(string)})
	}
	return tags
}

func parseOptionalExperience(val any) *WorkExperience {
	if val == nil {
		return nil
	}
	props := val.(neo4j.Node).Props
	return &WorkExperience{
		ID:        props["id"].(string),
		Summary:   props["summary"].(string),
		Company:   props["company"].(string),
		Title:     props["title"].(string),
		StartDate: props["startDate"].(string),
		EndDate:   props["endDate"].(string),
		Featured:  props["featured"].(bool),
	}
}

func safeString(val any) string {
	if s, ok := val.(string); ok {
		return s
	}
	return ""
}

// Helper to safely convert []any to []string
func toStringSlice(val any) []string {
	items, ok := val.([]any)
	if !ok {
		return nil
	}
	var result []string
	for _, item := range items {
		if s, ok := item.(string); ok {
			result = append(result, s)
		}
	}
	return result
}

func asString(record *neo4j.Record, key string) string {
	if val, ok := record.Get(key); ok && val != nil {
		if s, ok := val.(string); ok {
			return s
		}
	}
	return ""
}

func safeToStringSlice(record *neo4j.Record, key string) []string {
	val, ok := record.Get(key)
	if !ok || val == nil {
		return []string{}
	}
	return toStringSlice(val)
}
func asBool(record *neo4j.Record, key string) bool {
	if val, ok := record.Get(key); ok && val != nil {
		if b, ok := val.(bool); ok {
			return b
		}
	}
	return false
}
func toString(val interface{}) string {
	if s, ok := val.(string); ok {
		return s
	}
	return ""
}
func intFromInterface(val interface{}) int {
	switch v := val.(type) {
	case int:
		return v
	case int64:
		return int(v)
	case float64:
		return int(v)
	default:
		return 0
	}
}
func toBool(value any) bool {
	if b, ok := value.(bool); ok {
		return b
	}
	return false
}
func withReadSession(run func(tx neo4j.ManagedTransaction) (any, error)) (any, error) {
	ctx := context.Background()
	session := Neo4jDriver.NewSession(ctx, neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeRead,
	})
	defer session.Close(ctx)

	return session.ExecuteRead(ctx, run)
}
