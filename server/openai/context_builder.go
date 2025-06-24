package openai

import (
	"encoding/json"
	"fmt"
	"go-ai/db"
	"log"
	"strings"
)

// ContextIntent defines what the user is asking for
type ContextIntent struct {
	Skills     []string `json:"skills,omitempty"`
	Tags       []string `json:"tags,omitempty"`
	Projects   bool     `json:"projects"`
	Experience bool     `json:"experience"`
	Education  bool     `json:"education"`
	Hobbies    bool     `json:"hobbies"`
}

func extractKeywords(input string) []string {
	input = strings.ToLower(input)
	var matched []string
	seen := make(map[string]bool)

	synonyms := map[string][]string{
		"education":     {"school", "studied", "university"},
		"experience":    {"work", "job", "career", "employment"},
		"hobbies":       {"interests", "passions", "gaming", "music", "fun"},
		"graphql":       {"api", "query language"},
		"react":         {"frontend", "ui library"},
		"aws":           {"cloud", "deployment"},
		"docker":        {"container", "virtualization"},
		"cybersecurity": {"security", "infosec"},
		"hackathon":     {"competition", "coding challenge"},
		"python":        {"scripting", "language"},
	}

	known := make([]string, 0, len(synonyms))
	for k := range synonyms {
		known = append(known, k)
	}

	// check direct matches
	for _, kw := range known {
		if strings.Contains(input, kw) && !seen[kw] {
			matched = append(matched, kw)
			seen[kw] = true
			continue
		}
		for _, alt := range synonyms[kw] {
			if strings.Contains(input, alt) && !seen[kw] {
				matched = append(matched, kw)
				seen[kw] = true
				break
			}
		}
	}

	return matched
}
func BuildChatContext(userInput string) (string, error) {
	keywords := extractKeywords(userInput)

	if len(keywords) == 0 {
		// fallback to GPT-based intent classification
		intent, err := classifyUserIntent(userInput)
		if err != nil {
			return "", err
		}
		return buildContextFromIntent(intent)
	}

	return buildContextFromKeywords(keywords)
}

// BuildChatContext constructs a GPT-ready context string using local Neo4j queries based on extracted keywords
func buildContextFromKeywords(keywords []string) (string, error) {
	var sb strings.Builder

	for _, kw := range keywords {
		switch kw {
		case "education":
			edu, _ := db.GetEducation()
			for _, e := range edu {
				fmt.Fprintf(&sb, "🎓 %s — %s (%s–%s)\n", e.Degree, e.Institution, e.StartDate, e.EndDate)
			}
		case "experience", "work", "job":
			exp, _ := db.GetWorkExperience()
			for _, e := range exp {
				fmt.Fprintf(&sb, "💼 %s at %s (%s–%s)\n", e.Title, e.Company, e.StartDate, e.EndDate)
			}
		case "hobbies", "gaming", "music":
			hobbies, _ := db.GetHobbies()
			for _, h := range hobbies {
				if strings.Contains(strings.ToLower(h.Name), kw) || kw == "hobbies" {
					fmt.Fprintf(&sb, "🎮 %s: %s\n", h.Name, h.Description)
				}
			}
		default:
			// Match as skill or tag
			projectsBySkill, _ := db.FindProjectsBySkill(kw)
			for _, p := range projectsBySkill {
				fmt.Fprintf(&sb, "🛠️ Used skill [%s] in %s — %s\n", kw, p.Name, p.Description)
			}

			projectsByTag, _ := db.FindProjectsByTag(kw)
			for _, p := range projectsByTag {
				fmt.Fprintf(&sb, "🏷️ Tagged [%s]: %s — %s\n", kw, p.Name, p.Description)
			}
		}
	}

	if sb.Len() == 0 {
		return "🤷 I couldn’t find any relevant resume info.", nil
	}

	return sb.String(), nil
}

func buildContextFromIntent(intent ContextIntent) (string, error) {
	var sb strings.Builder

	if intent.Education {
		edu, _ := db.GetEducation()
		for _, e := range edu {
			fmt.Fprintf(&sb, "🎓 %s — %s (%s–%s)\n", e.Degree, e.Institution, e.StartDate, e.EndDate)
		}
	}

	if intent.Experience {
		exp, _ := db.GetWorkExperience()
		for _, e := range exp {
			fmt.Fprintf(&sb, "💼 %s at %s (%s–%s)\n", e.Title, e.Company, e.StartDate, e.EndDate)
		}
	}

	if intent.Hobbies {
		hobbies, _ := db.GetHobbies()
		for _, h := range hobbies {
			fmt.Fprintf(&sb, "🎮 %s: %s\n", h.Name, h.Description)
		}
	}

	for _, skill := range intent.Skills {
		projects, _ := db.FindProjectsBySkill(skill)
		for _, p := range projects {
			fmt.Fprintf(&sb, "🛠️ Used skill [%s] in %s — %s\n", skill, p.Name, p.Description)
		}
	}

	for _, tag := range intent.Tags {
		projects, _ := db.FindProjectsByTag(tag)
		for _, p := range projects {
			fmt.Fprintf(&sb, "🏷️ Tagged [%s]: %s — %s\n", tag, p.Name, p.Description)
		}
	}

	if intent.Projects && sb.Len() == 0 {
		projects, _ := db.GetAllProjects()
		for _, p := range projects {
			fmt.Fprintf(&sb, "📦 %s: %s\n", p.Name, p.Description)
		}
	}

	if sb.Len() == 0 {
		return "🤷 I couldn’t find any relevant resume info based on your message.", nil
	}

	return sb.String(), nil
}

func classifyUserIntent(userInput string) (ContextIntent, error) {
	log.Println("🧠 GPT intent classifier called (should only run if no keyword match)")
	prompt := fmt.Sprintf(`Given this message from a user:

"%s"

Determine what parts of a resume the user is asking about. Return only a JSON object like:

{
  "skills": ["React", "Docker", ...],
  "tags": ["Hackathon", "Cybersecurity"],
  "projects": true,
  "experience": true,
  "education": false,
  "hobbies": false
}

Use lowercase for all skill/tag names. Leave arrays empty if not mentioned.
`, userInput)

	resp, err := CallOpenAI([]db.ChatMessage{
		{Role: "system", Content: "You are a resume chatbot. Return a clean JSON object of what resume info the user is asking about."},
		{Role: "user", Content: prompt},
	}, "gpt-3.5-turbo")

	if err != nil {
		return ContextIntent{}, err
	}

	resp = strings.TrimSpace(resp)
	start := strings.Index(resp, "{")
	end := strings.LastIndex(resp, "}") + 1
	if start == -1 || end == -1 || end <= start {
		return ContextIntent{}, fmt.Errorf("⚠️ invalid JSON format returned:\n%s", resp)
	}

	var intent ContextIntent
	if err := json.Unmarshal([]byte(resp[start:end]), &intent); err != nil {
		return ContextIntent{}, fmt.Errorf("intent unmarshal failed: %w\n\nRaw response: %s", err, resp)
	}

	return intent, nil
}
