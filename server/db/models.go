package db

import "time"

type ChatEntry struct {
	UserID    string    `bson:"user_id"`
	Role      string    `bson:"role"`
	Content   string    `bson:"content"`
	Timestamp time.Time `bson:"timestamp"`
}
type Course struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Embedding   []float64 `json:"embedding"`
}
type Education struct {
	ID          string    `json:"id"`
	Institution string    `json:"institution"`
	Degree      string    `json:"degree"`
	Field       string    `json:"field"`
	Level       string    `json:"level"`
	StartDate   string    `json:"startDate"`
	EndDate     string    `json:"endDate"`
	Summary     string    `json:"summary"`
	Leadership  []string  `json:"leadership"`
	Embedding   []float64 `json:"embedding"`
}
type Hobby struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Embedding   []float64 `json:"embedding"`
}
type Person struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	Pronouns   string    `json:"pronouns"`
	Summary    string    `json:"summary"`
	Location   string    `json:"location"`
	BirthMonth string    `json:"birthMonth"`
	BirthYear  int64     `json:"birthYear"`
	VoiceTone  string    `json:"voiceTone"`
	Background []string  `json:"background"`
	Embedding  []float64 `json:"embedding"`
}

type Project struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	Type          string    `json:"type"`
	Institution   string    `json:"institution"`
	Description   string    `json:"description"`
	Contributions []string  `json:"contributions"`
	StartDate     string    `json:"startDate"`
	EndDate       string    `json:"endDate"`
	Featured      bool      `json:"featured"`
	Github        string    `json:"github"`
	Demo          string    `json:"demo"`
	Image         string    `json:"image"`
	Embedding     []float64 `json:"embedding"`
}
type Skill struct {
	Name string `json:"name"`
}
type Tag struct {
	Name string `json:"name"`
}
type Team struct {
	Name       string   `json:"name"`
	Role       string   `json:"role"`
	Period     string   `json:"period"`
	Highlights []string `json:"highlights"`
}
type WorkExperience struct {
	ID        string    `json:"id"`
	Company   string    `json:"company"`
	Title     string    `json:"title"`
	StartDate string    `json:"startDate"`
	EndDate   string    `json:"endDate"`
	Summary   string    `json:"summary"`
	Featured  bool      `json:"featured"`
	Embedding []float64 `json:"embedding"`
}
