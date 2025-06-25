package db

import "log"

type FilterClause struct {
	On       string `json:"on"`       // e.g., "Tag"
	Value    string `json:"value"`    // e.g., "Hyperpad"
	Relation string `json:"relation"` // e.g., "HAS_TAG"
}

// FindProjectsWithFilters dispatches filter-based queries for projects
func FindProjectsWithFilters(filters []FilterClause) ([]Project, error) {
	for _, f := range filters {
		switch {
		case f.On == "Tag":
			return FindProjectsByTag(f.Value)
		case f.On == "Skill":
			return FindProjectsBySkill(f.Value)
		case f.On == "Hobby":
			return FindProjectsByHobby(f.Value)
		case f.On == "Name":
			return SearchProjectsByName(f.Value)
		default:
			log.Printf("⚠️ Ignoring unsupported project filter: %+v\n", f)
		}
	}
	return GetAllProjectsSorted()
}

// FindWorkExperienceWithFilters supports filtering work experience by company
func FindWorkExperienceWithFilters(filters []FilterClause) ([]WorkExperience, error) {
	for _, f := range filters {
		switch f.On {
		case "Tag":
			return FindWorkExperienceByTag(f.Value)
		case "Company":
			return SearchWorkExperiencesByCompany(f.Value)
		case "Name":
			return SearchWorkExperiencesByName(f.Value)
		}
	}
	return GetAllWorkExperiences()
}

// FindEducationWithFilters filters education (future: by institution, field, etc.)
func FindEducationWithFilters(filters []FilterClause) ([]Education, error) {
	log.Println("in education")
	for _, f := range filters {
		if f.On == "Institution" {
			return SearchEducationByInstitution(f.Value)
		}
		if f.On == "Field" {
			return SearchEducationByField(f.Value)
		}
	}
	return GetAllEducationSorted()
}

// FindHobbiesWithFilters filters hobbies (future: by name or tag)
func FindHobbiesWithFilters(filters []FilterClause) ([]Hobby, error) {
	for _, f := range filters {
		if f.On == "Name" {
			return SearchHobbiesByName(f.Value)
		}
		if f.On == "Tag" {
			return SearchHobbiesByTag(f.Value)
		}
	}
	return GetAllHobbies()
}

// FindSkillsWithFilters filters skills (future: by tag or project)
func FindSkillsWithFilters(filters []FilterClause) ([]Skill, error) {
	for _, f := range filters {
		if f.On == "Name" {
			return SearchSkillsByName(f.Value)
		}
		if f.On == "Tag" {
			return SearchSkillsByTag(f.Value)
		}
	}
	return GetAllSkillsSorted()
}
