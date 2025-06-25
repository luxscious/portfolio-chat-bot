package db

import (
	"context"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

// ─────────────────────────────────────────────────────────────────────────────
// PUBLIC FUNCTIONS
// ─────────────────────────────────────────────────────────────────────────────

// GetAllWorkExperiencesSorted returns all WorkExperience nodes sorted by startDate.
func GetAllWorkExperiencesSorted() ([]WorkExperience, error) {
	query := `
		MATCH (w:WorkExperience)
		RETURN w
		ORDER BY w.startDate
	`
	return queryWorkExperiences(query, nil)
}

// SearchWorkExperiencesByCompany performs a case-insensitive search on company name.
func SearchWorkExperiencesByCompany(company string) ([]WorkExperience, error) {
	query := `
		MATCH (w:WorkExperience)
		WHERE toLower(w.company) CONTAINS toLower($company)
		RETURN w
		ORDER BY w.startDate
	`
	params := map[string]interface{}{"company": company}
	return queryWorkExperiences(query, params)
}

// SearchWorkExperiencesByName searches by company OR title.
func SearchWorkExperiencesByName(name string) ([]WorkExperience, error) {
	query := `
		MATCH (w:WorkExperience)
		WHERE toLower(w.company) CONTAINS toLower($name)
		   OR toLower(w.title) CONTAINS toLower($name)
		RETURN w
	`
	params := map[string]interface{}{"name": name}
	return queryWorkExperiences(query, params)
}

// FindWorkExperienceByTag returns work experiences associated with a given tag.
// (unchanged)
func FindWorkExperienceByTag(tag string) ([]WorkExperience, error) {
	query := `
		MATCH (w:WorkExperience)-[:HAS_TAG]->(t:Tag)
		WHERE toLower(t.name) = toLower($tag)
		RETURN w
	`
	params := map[string]interface{}{"tag": tag}
	return queryWorkExperiences(query, params)
}

// ListWorkExperienceCompanies extracts just the company names from all work experiences.
func ListWorkExperienceCompanies() ([]string, error) {
	work, err := GetAllWorkExperiencesSorted()
	if err != nil {
		return nil, err
	}
	var companies []string
	for _, w := range work {
		companies = append(companies, w.Company)
	}
	return companies, nil
}

// GetAllWorkExperiences returns all WorkExperience nodes without sorting.
func GetAllWorkExperiences() ([]WorkExperience, error) {
	query := `
		MATCH (w:WorkExperience)
		RETURN w
	`
	return queryWorkExperiences(query, nil)
}

// ─────────────────────────────────────────────────────────────────────────────
// INTERNAL HELPER
// ─────────────────────────────────────────────────────────────────────────────

// queryWorkExperiences executes a read transaction and maps WorkExperience nodes.
func queryWorkExperiences(cypher string, params map[string]interface{}) ([]WorkExperience, error) {
	ctx := context.Background()
	session := Neo4jDriver.NewSession(ctx, neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeRead,
	})
	defer session.Close(ctx)

	result, err := session.Run(ctx, cypher, params)
	if err != nil {
		return nil, err
	}

	var experiences []WorkExperience
	for result.Next(ctx) {
		record := result.Record()
		node, ok := record.Get("w")
		if !ok {
			continue
		}
		props := node.(neo4j.Node).Props
		experience := WorkExperience{
			ID:        toString(props["id"]),
			Title:     toString(props["title"]),
			Company:   toString(props["company"]),
			Summary:   toString(props["summary"]),
			StartDate: toString(props["startDate"]), // Ensure consistent naming
			EndDate:   toString(props["endDate"]),
			Featured:  toBool(props["featured"]),
		}
		experiences = append(experiences, experience)
	}

	if err = result.Err(); err != nil {
		return nil, err
	}

	return experiences, nil
}
