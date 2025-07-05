package openai

import (
	"fmt"
	"go-ai/db"
	"go-ai/ollama"
	"strings"
)

// BuildContextFromGraphPlan gathers relevant context from Neo4j based on a structured query plan.
func BuildContextFromGraphPlan(plan ollama.GraphQueryPlan) (string, error) {
	var contextParts []string

	// Filter out empty or placeholder values
	var validFilters []db.FilterClause
	for _, f := range plan.Filters {
		if f.Value != "" && f.Value != "null" {
			validFilters = append(validFilters, f)
		}
	}
	plan.Filters = validFilters

	for _, nodeType := range plan.TargetNodes {
		switch nodeType {
		case "Project":
			projects, err := db.FindProjectsWithFilters(plan.Filters)
			if err != nil || len(projects) == 0 {
				// Fallback: get all projects if none found
				projects, err = db.FindProjectsWithFilters(nil)
				if err != nil || len(projects) == 0 {
					continue
				}
			}
			var b strings.Builder
			b.WriteString("Relevant Projects:\n")
			for _, p := range projects {
				b.WriteString(fmt.Sprintf("- %s: %s\n", p.Name, p.Description))
				if len(p.Contributions) > 0 {
					b.WriteString("  Contributions:\n")
					for _, c := range p.Contributions {
						b.WriteString(fmt.Sprintf("    • %s\n", c))
					}
				}
			}
			contextParts = append(contextParts, b.String())

		case "WorkExperience":
			experiences, err := db.FindWorkExperienceWithFilters(plan.Filters)
			if err != nil || len(experiences) == 0 {
				// Fallback: get all work experiences if none found
				experiences, err = db.FindWorkExperienceWithFilters(nil)
				if err != nil || len(experiences) == 0 {
					continue
				}
			}
			var b strings.Builder
			b.WriteString("Work Experience:\n")
			for _, w := range experiences {
				b.WriteString(fmt.Sprintf("- %s at %s: %s\n", w.Title, w.Company, w.Summary))
			}
			contextParts = append(contextParts, b.String())

		case "Education":
			education, err := db.FindEducationWithFilters(plan.Filters)
			if err != nil || len(education) == 0 {
				// Fallback: get all education if none found
				education, err = db.FindEducationWithFilters(nil)
				if err != nil || len(education) == 0 {
					continue
				}
			}
			var b strings.Builder
			b.WriteString("Education:\n")
			for _, e := range education {
				b.WriteString(fmt.Sprintf("- %s at %s: %s\n", e.Degree, e.Institution, e.Summary))
			}
			contextParts = append(contextParts, b.String())

		case "Hobby":
			hobbies, err := db.FindHobbiesWithFilters(plan.Filters)
			if err != nil || len(hobbies) == 0 {
				// Fallback: get all hobbies if none found
				hobbies, err = db.FindHobbiesWithFilters(nil)
				if err != nil || len(hobbies) == 0 {
					continue
				}
			}
			var b strings.Builder
			b.WriteString("Hobbies:\n")
			for _, h := range hobbies {
				b.WriteString(fmt.Sprintf("- %s: %s\n", h.Name, h.Description))
			}
			contextParts = append(contextParts, b.String())

		case "Skill":
			skills, err := db.FindSkillsWithFilters(plan.Filters)
			if err != nil || len(skills) == 0 {
				// Fallback: get all skills if none found
				skills, err = db.FindSkillsWithFilters(nil)
				if err != nil || len(skills) == 0 {
					continue
				}
			}
			var b strings.Builder
			b.WriteString("Skills:\n")
			for _, s := range skills {
				b.WriteString(fmt.Sprintf("- %s\n", s.Name))
			}
			contextParts = append(contextParts, b.String())

		case "Person":
			person, err := db.GetPerson()
			if err != nil {
				continue
			}
			bio := fmt.Sprintf(
				"%s is a %s based in %s. With a background in %s and pronouns %s, she brings a love for problem solving and gaming into all her work. She’s especially passionate about cybersecurity, EV infrastructure, and full-stack development.",
				person.Name,
				strings.ToLower(person.Summary),
				person.Location,
				strings.Join(person.Background, " and "),
				person.Pronouns,
			)
			contextParts = append(contextParts, bio)
		}
	}

	return strings.Join(contextParts, "\n\n"), nil
}
