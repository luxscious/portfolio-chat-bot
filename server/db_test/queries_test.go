package db_test

import (
	"fmt"
	"go-ai/config"
	"go-ai/db"
	"testing"
)

func TestGetAllProjectsSorted(t *testing.T) {
	config.LoadEnv()
	db.InitNeo4j()
	defer db.Neo4jDriver.Close(nil)

	projects, err := db.GetAllProjectsSorted()
	if err != nil {
		t.Fatalf("❌ GetAllProjectsSorted failed: %v", err)
	}

	for _, p := range projects {
		fmt.Printf("📦 %s (%s) — %s\n", p.Name, p.ID, p.Description)
	}
}

func TestGetAllEducationSorted(t *testing.T) {
	config.LoadEnv()
	db.InitNeo4j()
	defer db.Neo4jDriver.Close(nil)

	education, err := db.GetAllEducationSorted()
	if err != nil {
		t.Fatalf("❌ GetAllEducationSorted failed: %v", err)
	}

	for _, e := range education {
		fmt.Printf("🎓 %s (%s) — %s\n", e.Degree, e.ID, e.Institution)
	}
}

func TestGetAllWorkExperiencesSorted(t *testing.T) {
	config.LoadEnv()
	db.InitNeo4j()
	defer db.Neo4jDriver.Close(nil)

	work, err := db.GetAllWorkExperiencesSorted()
	if err != nil {
		t.Fatalf("❌ GetAllWorkExperiencesSorted failed: %v", err)
	}

	for _, w := range work {
		fmt.Printf("💼 %s at %s (%s - %s)\n", w.Title, w.Company, w.StartDate, w.EndDate)
	}
}

func TestGetAllHobbies(t *testing.T) {
	config.LoadEnv()
	db.InitNeo4j()
	defer db.Neo4jDriver.Close(nil)

	hobbies, err := db.GetAllHobbies()
	if err != nil {
		t.Fatalf("❌ GetAllHobbies failed: %v", err)
	}

	for _, h := range hobbies {
		fmt.Printf("🎮 %s — %s\n", h.Name, h.Description)
	}
}

func TestGetAllSkillsSorted(t *testing.T) {
	config.LoadEnv()
	db.InitNeo4j()
	defer db.Neo4jDriver.Close(nil)

	skills, err := db.GetAllSkillsSorted()
	if err != nil {
		t.Fatalf("❌ GetAllSkillsSorted failed: %v", err)
	}

	for _, s := range skills {
		fmt.Printf("🛠️ %s\n", s.Name)
	}
}

func TestGetProjectDetails(t *testing.T) {
	config.LoadEnv()
	db.InitNeo4j()
	defer db.Neo4jDriver.Close(nil)

	projectID := "thesis-research"

	details, err := db.GetProjectDetails(projectID)
	if err != nil {
		t.Fatalf("❌ GetProjectDetails failed: %v", err)
	}

	fmt.Printf("🧠 Project: %s\n", details.Project.Name)
	fmt.Printf("📄 Description: %s\n", details.Project.Description)
	fmt.Printf("🏫 Institution: %s\n", details.Project.Institution)
	fmt.Printf("📅 Start: %s → End: %s\n", details.Project.StartDate, details.Project.EndDate)

	if details.Experience != nil {
		fmt.Printf("💼 Related Job: %s at %s\n", details.Experience.Title, details.Experience.Company)
	}

	fmt.Println("🔧 Skills Used:")
	for _, skill := range details.Skills {
		fmt.Printf("- %s\n", skill.Name)
	}

	fmt.Println("🏷️ Tags:")
	for _, tag := range details.Tags {
		fmt.Printf("- %s\n", tag.Name)
	}
}

func TestFindProjectsBySkill(t *testing.T) {
	config.LoadEnv()
	db.InitNeo4j()
	defer db.Neo4jDriver.Close(nil)

	skill := "React"
	projects, err := db.FindProjectsBySkill(skill)
	if err != nil {
		t.Fatalf("❌ FindProjectsBySkill failed: %v", err)
	}

	fmt.Printf("🛠️ Projects using skill: %s\n", skill)
	for _, p := range projects {
		fmt.Printf("📦 %s (%s): %s\n", p.Name, p.ID, p.Description)
	}
}

func TestFindProjectsByTag(t *testing.T) {
	config.LoadEnv()
	db.InitNeo4j()
	defer db.Neo4jDriver.Close(nil)

	tag := "Hackathon"
	projects, err := db.FindProjectsByTag(tag)
	if err != nil {
		t.Fatalf("❌ FindProjectsByTag failed: %v", err)
	}

	fmt.Printf("🏷️ Projects with tag: %s\n", tag)
	for _, p := range projects {
		fmt.Printf("📦 %s (%s): %s\n", p.Name, p.ID, p.Description)
	}
}

func TestSearchWorkExperiencesByCompany(t *testing.T) {
	config.LoadEnv()
	db.InitNeo4j()
	defer db.Neo4jDriver.Close(nil)

	companyName := "Hyperpad"
	experiences, err := db.SearchWorkExperiencesByCompany(companyName)
	if err != nil {
		t.Fatalf("❌ Failed to find experience for company %q: %v", companyName, err)
	}

	if len(experiences) == 0 {
		t.Fatalf("❌ No experiences found for company %q", companyName)
	}

	for _, exp := range experiences {
		fmt.Printf("✅ Found: %s at %s\nSummary: %s\n\n", exp.Title, exp.Company, exp.Summary)
	}
}

func TestFindProjectsByHobby(t *testing.T) {
	config.LoadEnv()
	db.InitNeo4j()
	defer db.Neo4jDriver.Close(nil)

	hobbyName := "Hackathons"
	projects, err := db.FindProjectsByHobby(hobbyName)
	if err != nil {
		t.Fatalf("❌ Error fetching projects for hobby %q: %v", hobbyName, err)
	}

	if len(projects) == 0 {
		t.Fatalf("❌ No projects found for hobby %q", hobbyName)
	}

	for _, p := range projects {
		fmt.Printf("✅ Project: %s\nDescription: %s\n\n", p.Name, p.Description)
	}
}

func TestFindTagsBySkill(t *testing.T) {
	config.LoadEnv()
	db.InitNeo4j()
	defer db.Neo4jDriver.Close(nil)

	skill := "React"
	tags, err := db.FindTagsBySkill(skill)
	if err != nil {
		t.Fatalf("❌ Error fetching tags for skill %q: %v", skill, err)
	}

	if len(tags) == 0 {
		t.Fatalf("❌ No tags found for skill %q", skill)
	}

	for _, tag := range tags {
		fmt.Printf("✅ Tag associated with skill %q: %s\n", skill, tag.Name)
	}
}

func TestSearchSkillsByTag(t *testing.T) {
	config.LoadEnv()
	db.InitNeo4j()
	defer db.Neo4jDriver.Close(nil)

	tag := "Frontend"
	skills, err := db.SearchSkillsByTag(tag)
	if err != nil {
		t.Fatalf("❌ Error fetching skills for tag %q: %v", tag, err)
	}

	if len(skills) == 0 {
		t.Fatalf("❌ No skills found for tag %q", tag)
	}

	for _, skill := range skills {
		fmt.Printf("✅ Skill used in %q projects: %s\n", tag, skill.Name)
	}
}
