package db

import (
	"context"
	"errors"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

// PROJECT-RELATED
// GetAllProjects retrieves all project nodes from Neo4j
func GetAllProjects() ([]Project, error) {
	session := Neo4jDriver.NewSession(context.Background(), neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeRead,
	})
	defer session.Close(context.Background())

	result, err := session.ExecuteRead(context.Background(), func(tx neo4j.ManagedTransaction) (any, error) {
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

		cypherResult, err := tx.Run(context.Background(), query, nil)
		if err != nil {
			return nil, err
		}

		var projects []Project

		for cypherResult.Next(context.Background()) {
			record := cypherResult.Record()

			id, _ := record.Get("id")
			name, _ := record.Get("name")
			description, _ := record.Get("description")
			institution, _ := record.Get("institution")
			image, _ := record.Get("image")
			featured, _ := record.Get("featured")
			contributions, _ := record.Get("contributions")
			startDate, _ := record.Get("startDate")
			endDate, _ := record.Get("endDate")
			demo, _ := record.Get("demo")
			github, _ := record.Get("github")

			project := Project{
				ID:            id.(string),
				Name:          name.(string),
				Description:   description.(string),
				Institution:   institution.(string),
				Image:         image.(string),
				Featured:      featured.(bool),
				Contributions: toStringSlice(contributions),
				StartDate:     startDate.(string),
				EndDate:       endDate.(string),
			}

			if demo != nil {
				project.Demo = demo.(string)
			}
			if github != nil {
				project.GitHub = github.(string)
			}
			projects = append(projects, project)
		}

		if err := cypherResult.Err(); err != nil {
			return nil, err
		}

		return projects, nil
	})
	if err != nil {
		return nil, err
	}

	// Now assert the type back to []models.Project
	return result.([]Project), nil
}

// // FindProjectsByTag returns projects associated with a specific tag (e.g., "React")
func FindProjectsByTag(tag string) ([]Project, error) {
	session := Neo4jDriver.NewSession(context.Background(), neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeRead,
	})
	defer session.Close(context.Background())

	result, err := session.ExecuteRead(context.Background(), func(tx neo4j.ManagedTransaction) (any, error) {
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

		params := map[string]any{"tag": tag}
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

// // FindProjectsBySkill returns projects that used a given skill
func FindProjectsBySkill(skill string) ([]Project, error) {
	session := Neo4jDriver.NewSession(context.Background(), neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeRead,
	})
	defer session.Close(context.Background())

	result, err := session.ExecuteRead(context.Background(), func(tx neo4j.ManagedTransaction) (any, error) {
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
				p.github AS github,
				p.type AS type
			ORDER BY p.startDate DESC
		`

		params := map[string]any{"skill": skill}
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

// // GetProjectDetails returns a project along with its skills, tags, and linked experience
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

			project := parseProjectNode(pNode)
			skills := parseSkillList(skillNodes)
			tags := parseTagList(tagNodes)
			experience := parseOptionalExperience(workNode)

			return ProjectDetails{
				Project:    project,
				Skills:     skills,
				Tags:       tags,
				Experience: experience,
			}, nil
		}

		return nil, errors.New("project not found")
	})

	if err != nil {
		return ProjectDetails{}, err
	}

	return result.(ProjectDetails), nil
}

// // EXPERIENCE-RELATED

// // GetExperienceTimeline returns all jobs ordered chronologically
func GetWorkExperience() ([]WorkExperience, error) {
	session := Neo4jDriver.NewSession(context.Background(), neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeRead,
	})
	defer session.Close(context.Background())

	result, err := session.ExecuteRead(context.Background(), func(tx neo4j.ManagedTransaction) (any, error) {
		query := `
			MATCH (w:WorkExperience)
			RETURN 
				w.id AS id,
				w.summary AS summary,
				w.company AS company,
				w.title AS title,
				w.startDate AS startDate,
				w.endDate AS endDate,
				w.featured AS featured
			ORDER BY w.startDate
		`

		res, err := tx.Run(context.Background(), query, nil)
		if err != nil {
			return nil, err
		}

		var workList []WorkExperience

		for res.Next(context.Background()) {
			record := res.Record()

			work := WorkExperience{
				ID:        asString(record, "id"),
				Summary:   asString(record, "summary"),
				Company:   asString(record, "company"),
				Title:     asString(record, "title"),
				StartDate: asString(record, "startDate"),
				EndDate:   asString(record, "endDate"),
				Featured:  asBool(record, "featured"),
			}

			workList = append(workList, work)
		}

		if err := res.Err(); err != nil {
			return nil, err
		}

		return workList, nil
	})

	if err != nil {
		return nil, err
	}

	return result.([]WorkExperience), nil
}

// // EDUCATION / HOBBIES / SKILLS

// // GetEducation returns all education entries
func GetEducation() ([]Education, error) {
	session := Neo4jDriver.NewSession(context.Background(), neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeRead,
	})
	defer session.Close(context.Background())

	result, err := session.ExecuteRead(context.Background(), func(tx neo4j.ManagedTransaction) (any, error) {
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

		res, err := tx.Run(context.Background(), query, nil)
		if err != nil {
			return nil, err
		}

		var educationList []Education

		for res.Next(context.Background()) {
			record := res.Record()

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

		if err := res.Err(); err != nil {
			return nil, err
		}

		return educationList, nil
	})

	if err != nil {
		return nil, err
	}

	return result.([]Education), nil
}

// // GetHobbies returns all hobby nodes
func GetHobbies() ([]Hobby, error) {
	session := Neo4jDriver.NewSession(context.Background(), neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeRead,
	})
	defer session.Close(context.Background())

	result, err := session.ExecuteRead(context.Background(), func(tx neo4j.ManagedTransaction) (any, error) {
		query := `
			MATCH (h:Hobby)
			RETURN 
				h.name AS name,
				h.description AS description
			ORDER BY h.name
		`

		res, err := tx.Run(context.Background(), query, nil)
		if err != nil {
			return nil, err
		}

		var hobbies []Hobby

		for res.Next(context.Background()) {
			record := res.Record()

			hobby := Hobby{
				Name:        asString(record, "name"),
				Description: asString(record, "description"),
			}

			hobbies = append(hobbies, hobby)
		}

		if err := res.Err(); err != nil {
			return nil, err
		}

		return hobbies, nil
	})

	if err != nil {
		return nil, err
	}

	return result.([]Hobby), nil
}

// // GetAllSkills returns all skill nodes
func GetAllSkills() ([]Skill, error) {
	session := Neo4jDriver.NewSession(context.Background(), neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeRead,
	})
	defer session.Close(context.Background())

	result, err := session.ExecuteRead(context.Background(), func(tx neo4j.ManagedTransaction) (any, error) {
		query := `
			MATCH (s:Skill)
			RETURN s.name AS name
			ORDER BY s.name
		`

		res, err := tx.Run(context.Background(), query, nil)
		if err != nil {
			return nil, err
		}

		var skills []Skill

		for res.Next(context.Background()) {
			record := res.Record()
			skills = append(skills, Skill{
				Name: asString(record, "name"),
			})
		}

		if err := res.Err(); err != nil {
			return nil, err
		}

		return skills, nil
	})

	if err != nil {
		return nil, err
	}

	return result.([]Skill), nil
}
