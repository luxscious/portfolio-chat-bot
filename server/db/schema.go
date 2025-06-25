package db

import (
	"context"
	"fmt"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

const GetAllNodeLabels = `
MATCH (n)
UNWIND labels(n) AS label
RETURN DISTINCT label
ORDER BY label
`

const GetAllSchemaRelationships = `
MATCH (a)-[r]->(b)
RETURN DISTINCT labels(a)[0] AS from, type(r) AS rel, labels(b)[0] AS to
ORDER BY from, rel, to
`

type GraphSchema struct {
	NodeLabels    []string
	Relationships []string
}

var CachedSchema GraphSchema

func LoadGraphSchemaOnce() error {
	session := Neo4jDriver.NewSession(context.Background(), neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeRead,
	})

	defer session.Close(context.Background())

	// Load Node Labels
	nodeLabels, err := session.ExecuteRead(context.Background(), func(tx neo4j.ManagedTransaction) (any, error) {
		res, err := tx.Run(context.Background(), GetAllNodeLabels, nil)
		if err != nil {
			return nil, err
		}
		var labels []string
		for res.Next(context.Background()) {
			label, _ := res.Record().Get("label")
			labels = append(labels, label.(string))
		}
		return labels, nil
	})
	if err != nil {
		return fmt.Errorf("failed to load node labels: %w", err)
	}

	// Load Relationships
	relationships, err := session.ExecuteRead(context.Background(), func(tx neo4j.ManagedTransaction) (any, error) {
		res, err := tx.Run(context.Background(), GetAllSchemaRelationships, nil)
		if err != nil {
			return nil, err
		}
		var rels []string
		for res.Next(context.Background()) {
			from, _ := res.Record().Get("from")
			rel, _ := res.Record().Get("rel")
			to, _ := res.Record().Get("to")
			rels = append(rels, fmt.Sprintf("(%s)-[:%s]->(%s)", from, rel, to))
		}
		return rels, nil
	})
	if err != nil {
		return fmt.Errorf("failed to load schema relationships: %w", err)
	}

	CachedSchema = GraphSchema{
		NodeLabels:    nodeLabels.([]string),
		Relationships: relationships.([]string),
	}
	return nil
}
