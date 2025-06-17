package main

import (
	"fmt"      //print to console or output to response
	"log"      //for logging messages or errors
	"net/http" //built-in HTTP server
)

func main() {
	// Health check endpoint
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "👋 Resume Chatbot backend is running!")
	})

	// Future endpoint: POST /chat
	http.HandleFunc("/chat", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
			return
		}
		fmt.Fprintln(w, "This is where the chatbot response will go.")
	})

	port := ":8080"
	log.Printf("✅ Server started on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
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
