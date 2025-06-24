package db

import "time"

type ChatEntry struct {
	UserID    string    `bson:"user_id"`
	Role      string    `bson:"role"`
	Content   string    `bson:"content"`
	Timestamp time.Time `bson:"timestamp"`
}

// 🎓 EDUCATION
type Education struct {
	ID          string   `json:"id"`
	Summary     string   `json:"summary"`
	Institution string   `json:"institution"`
	Field       string   `json:"field"`
	EndDate     string   `json:"endDate"`
	Level       string   `json:"level"`
	Degree      string   `json:"degree"`
	StartDate   string   `json:"startDate"`
	Leadership  []string `json:"leadership,omitempty"`
}

// 📚 COURSE
type Course struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

// 💼 WORK EXPERIENCE
type WorkExperience struct {
	ID        string `json:"id"`
	Summary   string `json:"summary"`
	Company   string `json:"company"`
	Title     string `json:"title"`
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
	Featured  bool   `json:"featured"`
}

// 🧠 HOBBY
type Hobby struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// 🚀 PROJECT
type Project struct {
	ID            string   `json:"id"`
	Name          string   `json:"name"`
	Description   string   `json:"description"`
	Institution   string   `json:"institution"`
	Image         string   `json:"image"`
	Featured      bool     `json:"featured"`
	Contributions []string `json:"contributions"`
	EndDate       string   `json:"endDate"`
	StartDate     string   `json:"startDate"`
	Demo          string   `json:"demo,omitempty"`
	GitHub        string   `json:"github,omitempty"`
}

// 🧩 SKILL
type Skill struct {
	Name string `json:"name"`
}

// 🙋‍♀️ PERSON
type Person struct {
	ID         string   `json:"id"`
	Name       string   `json:"name"`
	Summary    string   `json:"summary"`
	BirthMonth string   `json:"birthMonth"`
	BirthYear  int      `json:"birthYear"`
	Background []string `json:"background"`
	VoiceTone  string   `json:"voiceTone"`
	Location   string   `json:"location"`
	Pronouns   string   `json:"pronouns"`
}

type ProjectDetails struct {
	Project    Project         `json:"project"`
	Skills     []Skill         `json:"skills,omitempty"`
	Tags       []Tag           `json:"tags,omitempty"`
	Experience *WorkExperience `json:"experience,omitempty"`
}
type Tag struct {
	Name string `json:"name"`
}
