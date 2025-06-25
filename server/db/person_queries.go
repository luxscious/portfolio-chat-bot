package db

import (
	"context"
	"fmt"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func GetPerson() (*Person, error) {
	ctx := context.Background()
	session := Neo4jDriver.NewSession(ctx, neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeRead,
	})
	defer session.Close(ctx)

	query := `MATCH (p:Person) RETURN p LIMIT 1`
	result, err := session.Run(ctx, query, nil)
	if err != nil {
		return nil, err
	}

	if result.Next(ctx) {
		record := result.Record()
		node, ok := record.Get("p")
		if !ok {
			return nil, fmt.Errorf("person node not found")
		}

		props := node.(neo4j.Node).Props
		person := &Person{
			ID:         toString(props["id"]),
			Name:       toString(props["name"]),
			Summary:    toString(props["summary"]),
			Pronouns:   toString(props["pronouns"]),
			Location:   toString(props["location"]),
			BirthMonth: toString(props["birthMonth"]),
			BirthYear:  intFromInterface(props["birthYear"]),
		}

		if bg, ok := props["background"].([]any); ok {
			for _, val := range bg {
				person.Background = append(person.Background, toString(val))
			}
		}

		if vt, ok := props["voiceTone"].(string); ok {
			person.VoiceTone = vt
		}

		return person, nil
	}

	return nil, fmt.Errorf("no person node found")
}
