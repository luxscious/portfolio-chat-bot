package db

import (
	"context"
	"errors"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

// ─────────────────────────────────────────────────────────────────────────────
// PUBLIC FUNCTIONS
// ─────────────────────────────────────────────────────────────────────────────

// SearchProjectsByName returns projects where the name matches input (case-insensitive).
// Previously: FindProjectsByName
func SearchProjectsByName(name string) ([]Project, error) {
	query := `
		MATCH (p:Project)
		WHERE toLower(p.name) CONTAINS toLower($name)
		RETURN p
	`
	params := map[string]interface{}{"name": name}
	return queryProjects(query, params)
}

// GetAllProjectsSorted returns all projects sorted by start date.
// Previously: GetAllProjects
func GetAllProjectsSorted() ([]Project, error) {
	query := `
		MATCH (p:Project)
		RETURN 
			p.id AS id,
			p.name AS name,
			p.description AS description,
			p.institution AS institution,
			p.image AS image,
			p.featured AS featured,
			p.contributions AS contributions,
			p.startDate AS startDate,
			p.endDate AS endDate,
			p.demo AS demo,
			p.github AS github
		ORDER BY p.startDate DESC
	`
	return runProjectResultQuery(query, nil)
}

// ListProjectNames returns all project names only.
// Previously: GetAllProjectNames
func ListProjectNames() ([]string, error) {
	projects, err := GetAllProjectsSorted()
	if err != nil {
		return nil, err
	}
	var names []string
	for _, p := range projects {
		names = append(names, p.Name)
	}
	return names, nil
}

// FindProjectsByTag returns projects associated with a specific tag.
func FindProjectsByTag(tag string) ([]Project, error) {
	query := `
		MATCH (p:Project)-[:HAS_TAG]->(t:Tag {name: $tag})
		RETURN 
			p.id AS id,
			p.name AS name,
			p.description AS description,
			p.institution AS institution,
			p.image AS image,
			p.featured AS featured,
			p.contributions AS contributions,
			p.startDate AS startDate,
			p.endDate AS endDate,
			p.demo AS demo,
			p.github AS github
		ORDER BY p.startDate DESC
	`
	return runProjectResultQuery(query, map[string]any{"tag": tag})
}

// FindProjectsBySkill returns projects that use a specific skill.
func FindProjectsBySkill(skill string) ([]Project, error) {
	query := `
		MATCH (p:Project)-[:USES]->(s:Skill {name: $skill})
		RETURN 
			p.id AS id,
			p.name AS name,
			p.description AS description,
			p.institution AS institution,
			p.image AS image,
			p.featured AS featured,
			p.contributions AS contributions,
			p.startDate AS startDate,
			p.endDate AS endDate,
			p.demo AS demo,
			p.github AS github
		ORDER BY p.startDate DESC
	`
	return runProjectResultQuery(query, map[string]any{"skill": skill})
}

// FindProjectsConnectedToHobby returns projects linked to a specific hobby.
// Previously: FindProjectsConnectedToHobby
func FindProjectsByHobby(hobbyName string) ([]Project, error) {
	query := `
		MATCH (h:Hobby {name: $hobbyName})-[:INSPIRED]->(p:Project)
		RETURN 
			p.id AS id,
			p.name AS name,
			p.description AS description,
			p.institution AS institution,
			p.image AS image,
			p.featured AS featured,
			p.contributions AS contributions,
			p.startDate AS startDate,
			p.endDate AS endDate,
			p.demo AS demo,
			p.github AS github
		ORDER BY p.startDate DESC
	`
	return runProjectResultQuery(query, map[string]any{"hobbyName": hobbyName})
}

// GetProjectDetails returns a single project with its connected skills, tags, and work experience.
func GetProjectDetails(projectID string) (ProjectDetails, error) {
	session := Neo4jDriver.NewSession(context.Background(), neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeRead,
	})
	defer session.Close(context.Background())

	result, err := session.ExecuteRead(context.Background(), func(tx neo4j.ManagedTransaction) (any, error) {
		query := `
			MATCH (p:Project {id: $projectID})
			OPTIONAL MATCH (p)-[:USES]->(s:Skill)
			OPTIONAL MATCH (p)-[:HAS_TAG]->(t:Tag)
			OPTIONAL MATCH (p)-[:WORKED_ON]->(w:WorkExperience)
			RETURN 
				p,
				collect(DISTINCT s) AS skills,
				collect(DISTINCT t) AS tags,
				w
		`
		params := map[string]any{"projectID": projectID}
		res, err := tx.Run(context.Background(), query, params)
		if err != nil {
			return nil, err
		}

		if res.Next(context.Background()) {
			record := res.Record()
			pNode, _ := record.Get("p")
			skillNodes, _ := record.Get("skills")
			tagNodes, _ := record.Get("tags")
			workNode, _ := record.Get("w")

			return ProjectDetails{
				Project:    parseProjectNode(pNode),
				Skills:     parseSkillList(skillNodes),
				Tags:       parseTagList(tagNodes),
				Experience: parseOptionalExperience(workNode),
			}, nil
		}

		return nil, errors.New("project not found")
	})

	if err != nil {
		return ProjectDetails{}, err
	}

	return result.(ProjectDetails), nil
}

// ─────────────────────────────────────────────────────────────────────────────
// INTERNAL HELPERS
// ─────────────────────────────────────────────────────────────────────────────

// queryProjects is a fallback lightweight query that returns minimal Project data.
// Previously: runProjectQuery
func queryProjects(cypher string, params map[string]interface{}) ([]Project, error) {
	ctx := context.Background()
	session := Neo4jDriver.NewSession(ctx, neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeRead,
	})
	defer session.Close(ctx)

	result, err := session.Run(ctx, cypher, params)
	if err != nil {
		return nil, err
	}

	var projects []Project
	for result.Next(ctx) {
		record := result.Record()
		node, ok := record.Get("p")
		if !ok {
			continue
		}
		props := node.(neo4j.Node).Props
		project := Project{
			ID:            toString(props["id"]),
			Name:          toString(props["name"]),
			Description:   toString(props["description"]),
			Image:         toString(props["image"]),
			Demo:          toString(props["demo"]),
			Contributions: toStringSlice(props["contributions"]),
		}
		projects = append(projects, project)
	}

	if err = result.Err(); err != nil {
		return nil, err
	}

	return projects, nil
}

// runProjectResultQuery is the full record-based Cypher processor
func runProjectResultQuery(query string, params map[string]interface{}) ([]Project, error) {
	session := Neo4jDriver.NewSession(context.Background(), neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeRead,
	})
	defer session.Close(context.Background())

	result, err := session.ExecuteRead(context.Background(), func(tx neo4j.ManagedTransaction) (any, error) {
		res, err := tx.Run(context.Background(), query, params)
		if err != nil {
			return nil, err
		}

		var projects []Project
		for res.Next(context.Background()) {
			record := res.Record()
			project := Project{
				ID:            asString(record, "id"),
				Name:          asString(record, "name"),
				Description:   asString(record, "description"),
				Institution:   asString(record, "institution"),
				Image:         asString(record, "image"),
				Featured:      asBool(record, "featured"),
				Contributions: safeToStringSlice(record, "contributions"),
				StartDate:     asString(record, "startDate"),
				EndDate:       asString(record, "endDate"),
				Demo:          asString(record, "demo"),
				GitHub:        asString(record, "github"),
			}
			projects = append(projects, project)
		}

		if err := res.Err(); err != nil {
			return nil, err
		}
		return projects, nil
	})

	if err != nil {
		return nil, err
	}
	return result.([]Project), nil
}
