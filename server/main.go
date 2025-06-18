package main

import (
	"bytes"
	"fmt" //print to console or output to response
	"go-ai/config"
	"go-ai/db"
	"log"      //for logging messages or errors
	"net/http" //built-in HTTP server
	"text/template"
)

var systemPrompt string
// Helper Functions
func buildSystemPrompt(p PersonaContext) (string, error) {
	tmpl, err := template.New("prompt").Parse(p.PromptTemplate)
	if err != nil {
		return "", fmt.Errorf("template parsing failed: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, p); err != nil {
		return "", fmt.Errorf("template execution failed: %w", err)
	}

	return buf.String(), nil
}
func stringifyResume(resume *ResumeData) string {
	text := resume.PersonaContext.Summary + "\n\n"

	for _, exp := range resume.WorkExperience {
		text += fmt.Sprintf("Work: %s – %s\n", exp.Name, exp.Description)
	}

	for _, edu := range resume.Education {
		text += fmt.Sprintf("Education: %s – %s\n", edu.Name, edu.Description)
	}

	for _, proj := range resume.Projects {
		text += fmt.Sprintf("Project: %s – %s\n", proj.Name, proj.Description)
	}

	return text
}

// init our env variables before main is called 
func init() {
	config.LoadEnv()
	db.InitMongo()
	
	resume, err := loadResume("resume.json")
	if err != nil {
		log.Fatalf("Failed to load resume: %v", err)
	}

	fmt.Println("✅ System prompt initialized")
	fmt.Println("System prompt preview:")
	systemPrompt, err = buildSystemPrompt(resume.PersonaContext)
	if err != nil {
		log.Fatalf("Failed to build system prompt: %v", err)
	}	
	fmt.Println("System prompt preview:")
	fmt.Println(systemPrompt)
	fmt.Printf("🔎 %d education entries, %d work entries, %d projects loaded\n",
	len(resume.Education), len(resume.WorkExperience), len(resume.Projects))
	
	
	
	// Call embedding function
	resumeText := stringifyResume(resume)
	vector, err := generateEmbedding(resumeText)
	if err != nil {
		log.Fatalf("Embedding failed: %v", err)
	}
	fmt.Printf("✅ Embedding generated! First 5 values: %v\n", vector[:5])
	
}
func main() {

	port := config.GetServerPort()
	log.Printf("✅ Server started on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, RegisterRoutes()))
}










// import (
// 	"encoding/json"
// 	"fmt"
// 	"io/ioutil"
// 	"log"
// 	"os"
// 	"strings"
// )

// type Identity struct {
// 	Name      string `json:"name"`
// 	Birthdate string `json:"birthdate"`
// 	Pronouns  string `json:"pronouns"`
// 	Background string `json:"background"`
// 	Location  string `json:"location"`
// }

// type PersonaContext struct {
// 	Summary   string   `json:"summary"`
// 	VoiceTone string   `json:"voiceTone"`
// 	Identity  Identity `json:"identity"`
// 	Values    []string `json:"values"`
// }

// type ResumeData struct {
// 	PersonaContext PersonaContext `json:"personaContext"`
// 	// You can add Education, Projects, etc. here later if needed
// }

// func loadResumeData(path string) (*ResumeData, error) {
// 	file, err := os.Open(path)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer file.Close()

// 	bytes, err := ioutil.ReadAll(file)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var resume ResumeData
// 	err = json.Unmarshal(bytes, &resume)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &resume, nil
// }


// func buildSystemPrompt(p PersonaContext) string {
// 	return fmt.Sprintf(`
// You are %s, a %s software engineer and cybersecurity researcher based in %s.
// Use first-person language. Your tone is: %s

// Summary:
// %s

// You value:
// - %s

// You can answer questions about your education, experience, projects, skills, interests, and approach to problem solving. Keep it human, helpful, and in your voice.
// 	`,
// 		p.Identity.Name,
// 		p.Identity.Background,
// 		p.Identity.Location,
// 		p.VoiceTone,
// 		p.Summary,
// 		strings.Join(p.Values, "\n- "),
// 	)
// }


// func main() {
// 	resume, err := loadResumeData("resume.json")
// 	if err != nil {
// 		log.Fatalf("Failed to load resume.json: %v", err)
// 	}

// 	systemPrompt := buildSystemPrompt(resume.PersonaContext)
// 	fmt.Println("Generated system prompt:\n", systemPrompt)

// 	// ⬇️ You’ll eventually pass this into your OpenAI call
// }
