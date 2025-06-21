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

func FlattenResume(data *ResumeData) string {
	var out strings.Builder

	// Persona
	out.WriteString("Persona:\n")
	out.WriteString(fmt.Sprintf("Summary: %s\n", data.PersonaContext.Summary))
	out.WriteString(fmt.Sprintf("Voice/Tone: %s\n", data.PersonaContext.VoiceTone))
	out.WriteString(fmt.Sprintf("Location: %s\n", data.PersonaContext.Identity.Location))
	out.WriteString(fmt.Sprintf("Values: %s\n", strings.Join(data.PersonaContext.Values, ", ")))
	out.WriteString("\n\n")

	// Education
	for _, edu := range data.Education {
		out.WriteString(fmt.Sprintf("Education: %s at %s (%s to %s)\n", edu.Name, edu.Institution, edu.StartDate, edu.EndDate))
		out.WriteString(fmt.Sprintf("Description: %s\n", edu.Description))
		if len(edu.Skills) > 0 {
			out.WriteString(fmt.Sprintf("Skills: %s\n", strings.Join(edu.Skills, ", ")))
		}
		if len(edu.Leadership) > 0 {
			out.WriteString(fmt.Sprintf("Leadership: %s\n", strings.Join(edu.Leadership, ", ")))
		}
		out.WriteString("\n")
	}

	// Work
	for _, exp := range data.WorkExperience {
		out.WriteString(fmt.Sprintf("Experience: %s at %s (%s to %s)\n", exp.Name, exp.Institution, exp.StartDate, exp.EndDate))
		out.WriteString(fmt.Sprintf("Description: %s\n", exp.Description))
		if len(exp.Skills) > 0 {
			out.WriteString(fmt.Sprintf("Skills: %s\n", strings.Join(exp.Skills, ", ")))
		}
		for _, proj := range exp.Projects {
			out.WriteString(fmt.Sprintf("Project: %s — %s\n", proj.Title, proj.Description))
		}
		out.WriteString("\n")
	}

	// Projects
	for _, proj := range data.Projects {
		out.WriteString(fmt.Sprintf("Project: %s (%s to %s)\n", proj.Name, proj.StartDate, proj.EndDate))
		out.WriteString(fmt.Sprintf("Description: %s\n", proj.Description))
		if len(proj.Skills) > 0 {
			out.WriteString(fmt.Sprintf("Skills: %s\n", strings.Join(proj.Skills, ", ")))
		}
		out.WriteString(fmt.Sprintf("Tags: %s\n", strings.Join(proj.Tags, ", ")))
		out.WriteString("\n")
	}

	// Hobbies
	for _, h := range data.Hobbies {
		out.WriteString(fmt.Sprintf("Hobby: %s — %s\n", h.Name, h.Description))
		if len(h.Highlights) > 0 {
			out.WriteString(fmt.Sprintf("Highlights: %s\n", strings.Join(h.Highlights, ", ")))
		}
		out.WriteString("\n")
	}

	return out.String()
}