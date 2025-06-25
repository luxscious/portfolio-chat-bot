package db

import (
	"context"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

// ─────────────────────────────────────────────────────────────────────────────
// PUBLIC FUNCTIONS
// ─────────────────────────────────────────────────────────────────────────────

// GetAllEducationSorted returns all education entries sorted by start date.
func GetAllEducationSorted() ([]Education, error) {
	query := `
		MATCH (e:Education)
		RETURN 
			e.id AS id,
			e.summary AS summary,
			e.institution AS institution,
			e.field AS field,
			e.degree AS degree,
			e.level AS level,
			e.startDate AS startDate,
			e.endDate AS endDate,
			COALESCE(e.leadership, []) AS leadership
		ORDER BY e.startDate
	`

	ctx := context.Background()
	session := Neo4jDriver.NewSession(ctx, neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeRead,
	})
	defer session.Close(ctx)

	result, err := session.Run(ctx, query, nil)
	if err != nil {
		return nil, err
	}

	var educationList []Education
	for result.Next(ctx) {
		record := result.Record()

		edu := Education{
			ID:          asString(record, "id"),
			Summary:     asString(record, "summary"),
			Institution: asString(record, "institution"),
			Field:       asString(record, "field"),
			Degree:      asString(record, "degree"),
			Level:       asString(record, "level"),
			StartDate:   asString(record, "startDate"),
			EndDate:     asString(record, "endDate"),
			Leadership:  safeToStringSlice(record, "leadership"),
		}
		educationList = append(educationList, edu)
	}

	if err = result.Err(); err != nil {
		return nil, err
	}

	return educationList, nil
}

// SearchEducationByInstitution returns education nodes matching the institution.
func SearchEducationByInstitution(institution string) ([]Education, error) {
	query := `
		MATCH (e:Education)
		WHERE toLower(e.institution) CONTAINS toLower($institution)
		RETURN e
	`
	params := map[string]interface{}{"institution": institution}
	return queryEducations(query, params)
}

// SearchEducationByField returns education nodes matching the field.
func SearchEducationByField(field string) ([]Education, error) {
	query := `
		MATCH (e:Education)
		WHERE toLower(e.field) CONTAINS toLower($field)
		RETURN e
	`
	params := map[string]interface{}{"field": field}
	return queryEducations(query, params)
}

// ─────────────────────────────────────────────────────────────────────────────
// INTERNAL HELPER
// ─────────────────────────────────────────────────────────────────────────────

// queryEducations executes a generic query and returns basic education data.
func queryEducations(cypher string, params map[string]interface{}) ([]Education, error) {
	ctx := context.Background()
	session := Neo4jDriver.NewSession(ctx, neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeRead,
	})
	defer session.Close(ctx)

	result, err := session.Run(ctx, cypher, params)
	if err != nil {
		return nil, err
	}

	var educations []Education
	for result.Next(ctx) {
		record := result.Record()
		node, ok := record.Get("e")
		if !ok {
			continue
		}

		props := node.(neo4j.Node).Props
		edu := Education{
			ID:          toString(props["id"]),
			Degree:      toString(props["degree"]),
			Institution: toString(props["institution"]),
			Field:       toString(props["field"]),
			Summary:     toString(props["summary"]),
		}
		educations = append(educations, edu)
	}

	if err = result.Err(); err != nil {
		return nil, err
	}

	return educations, nil
}
