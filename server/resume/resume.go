package resume

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

// Identity holds personal identity details
type Identity struct {
	Name       string `json:"name"`
	Birthdate  string `json:"birthdate"`
	Pronouns   string `json:"pronouns"`
	Background string `json:"background"`
	Location   string `json:"location"`
}

// PersonaContext holds personal summary and tone
type PersonaContext struct {
	Summary   string   `json:"summary"`
	VoiceTone string   `json:"voiceTone"`
	Identity  Identity `json:"identity"`
	Values    []string `json:"values"`
	PromptTemplate string   `json:"promptTemplate"`
}

// Experience block (work, teaching)
type Experience struct {
	ID          string   `json:"id"`
	StartDate   string   `json:"startDate"`
	EndDate     string   `json:"endDate"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Institution string   `json:"institution"`
	Projects    []struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	} `json:"projects,omitempty"`
	Skills   []string `json:"skills"`
	Tags     []string `json:"tags"`
	Demo     string   `json:"demo,omitempty"`
	Type     string   `json:"type"`
	Featured bool     `json:"featured"`
}

// Project block (independent or hackathons)
type Project struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	StartDate   string   `json:"startDate"`
	EndDate     string   `json:"endDate"`
	Description string   `json:"description"`
	Institution string   `json:"institution"`
	Skills      []string `json:"skills"`
	Tags        []string `json:"tags"`
	Demo        string   `json:"demo,omitempty"`
	Image       string   `json:"image,omitempty"`
	Type        string   `json:"type"`
	Featured    bool     `json:"featured"`
	Github      string   `json:"github,omitempty"`
	Placement string `json:"placement,omitempty"` // Only used for hackathons
}

// Education block
type Education struct {
	ID          string   `json:"id"`
	StartDate   string   `json:"startDate"`
	EndDate     string   `json:"endDate"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Institution string   `json:"institution"`
	Skills      []string `json:"skills"`
	Leadership  []string `json:"leadership,omitempty"`
	Tags        []string `json:"tags"`
	Type        string   `json:"type"`
	Featured    bool     `json:"featured"`
}
type Hobby struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Highlights  []string `json:"highlights,omitempty"`
	Tags        []string `json:"tags"`
}

// ResumeData is the full top-level structure
type ResumeData struct {
	PersonaContext PersonaContext `json:"personaContext"`
	Education      []Education    `json:"education"`
	WorkExperience []Experience   `json:"work_experience"`
	Projects       []Project      `json:"projects"`
	Hobbies        []Hobby       `json:"hobbies"`
}


// loadResume loads and parses resume.json
func LoadResume(path string) (*ResumeData, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var resume ResumeData
	if err := json.Unmarshal(byteValue, &resume); err != nil {
		return nil, err
	}

	fmt.Println("✅ Resume loaded for:", resume.PersonaContext.Identity.Name)
	return &resume, nil
}

func (r *ResumeData) FlattenResume() string {
	var sb strings.Builder

	sb.WriteString("Persona Summary:\n")
	sb.WriteString(r.PersonaContext.Summary + "\n\n")

	sb.WriteString("Voice & Values:\n")
	sb.WriteString("Tone: " + r.PersonaContext.VoiceTone + "\n")
	sb.WriteString("Values: " + strings.Join(r.PersonaContext.Values, ", ") + "\n\n")

	sb.WriteString("Education:\n")
	for _, edu := range r.Education {
		sb.WriteString(fmt.Sprintf("- %s at %s (%s – %s)\n  %s\n",
			edu.Name, edu.Institution, edu.StartDate, edu.EndDate, edu.Description))
	}
	sb.WriteString("\n")

	sb.WriteString("Work Experience:\n")
	for _, exp := range r.WorkExperience {
		sb.WriteString(fmt.Sprintf("- %s at %s (%s – %s)\n  %s\n",
			exp.Name, exp.Institution, exp.StartDate, exp.EndDate, exp.Description))
		for _, proj := range exp.Projects {
			sb.WriteString(fmt.Sprintf("    · Project: %s — %s\n", proj.Title, proj.Description))
		}
	}
	sb.WriteString("\n")

	sb.WriteString("Projects:\n")
	for _, proj := range r.Projects {
		sb.WriteString(fmt.Sprintf("- %s (%s – %s)\n  %s\n",
			proj.Name, proj.StartDate, proj.EndDate, proj.Description))
	}

	sb.WriteString("\nHobbies:\n")
	for _, h := range r.Hobbies {
		sb.WriteString(fmt.Sprintf("- %s: %s\n", h.Name, h.Description))
	}

	return sb.String()
}