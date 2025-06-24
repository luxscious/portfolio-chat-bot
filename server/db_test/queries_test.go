package db_test

import (
	"fmt"
	"go-ai/config"
	"go-ai/db"
	"testing"
)

func TestGetAllProjects(t *testing.T) {
	config.LoadEnv()
	db.InitNeo4j()
	defer db.Neo4jDriver.Close(nil)

	projects, err := db.GetAllProjects()
	if err != nil {
		t.Fatalf("❌ GetAllProjects failed: %v", err)
	}

	for _, p := range projects {
		fmt.Printf("📦 %s (%s) — %s\n", p.Name, p.ID, p.Description)
	}
}

func TestGetEducation(t *testing.T) {
	config.LoadEnv()
	db.InitNeo4j()
	defer db.Neo4jDriver.Close(nil)

	education, err := db.GetEducation()
	if err != nil {
		t.Fatalf("❌ GetEducation failed: %v", err)
	}

	for _, e := range education {
		fmt.Printf("🎓 %s (%s) — %s\n", e.Degree, e.ID, e.Institution)
	}
}

func TestGetWorkExperience(t *testing.T) {
	config.LoadEnv()
	db.InitNeo4j()
	defer db.Neo4jDriver.Close(nil)

	work, err := db.GetWorkExperience()
	if err != nil {
		t.Fatalf("❌ GetWorkExperience failed: %v", err)
	}

	for _, w := range work {
		fmt.Printf("💼 %s at %s (%s - %s)\n", w.Title, w.Company, w.StartDate, w.EndDate)
	}
}

func TestGetHobbies(t *testing.T) {
	config.LoadEnv()
	db.InitNeo4j()
	defer db.Neo4jDriver.Close(nil)

	hobbies, err := db.GetHobbies()
	if err != nil {
		t.Fatalf("❌ GetHobbies failed: %v", err)
	}

	for _, h := range hobbies {
		fmt.Printf("🎮 %s — %s\n", h.Name, h.Description)
	}
}

func TestGetAllSkills(t *testing.T) {
	config.LoadEnv()
	db.InitNeo4j()
	defer db.Neo4jDriver.Close(nil)

	skills, err := db.GetAllSkills()
	if err != nil {
		t.Fatalf("❌ GetAllSkills failed: %v", err)
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

	skill := "React" // Replace with a real skill in your DB
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

	tag := "Hackathon" // Replace with a real tag in your DB
	projects, err := db.FindProjectsByTag(tag)
	if err != nil {
		t.Fatalf("❌ FindProjectsByTag failed: %v", err)
	}

	fmt.Printf("🏷️ Projects with tag: %s\n", tag)
	for _, p := range projects {
		fmt.Printf("📦 %s (%s): %s\n", p.Name, p.ID, p.Description)
	}
}
