package main

import (
	"bytes"
	"fmt" //print to console or output to response
	"go-ai/config"
	"go-ai/db"
	"go-ai/openai"
	"go-ai/resume"
	"log"      //for logging messages or errors
	"net/http" //built-in HTTP server
	"text/template"
)

var systemPrompt string
// Helper Functions
func buildSystemPrompt(p resume.PersonaContext) (string, error) {
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

var global_resume_Data  *resume.ResumeData 
// init our env variables before main is called 
func init() {
	config.LoadEnv()
	db.InitMongo()

	// Load and parse resume JSON
	resumeData, err := resume.LoadResume("./resume/resume.json")
	if err != nil {
		log.Fatalf("Failed to load resume: %v", err)
	}
	global_resume_Data = resumeData

	// Build system prompt
	systemPrompt, err = buildSystemPrompt(resumeData.PersonaContext)
	if err != nil {
		log.Fatalf("Failed to build system prompt: %v", err)
	}
	fmt.Println("✅ System prompt initialized")
	fmt.Println("System prompt preview:")
	fmt.Println(systemPrompt)

	// Print resume stats
	fmt.Printf("🔎 %d education entries, %d work entries, %d projects loaded\n",
		len(resumeData.Education), len(resumeData.WorkExperience), len(resumeData.Projects))

	// Perform RAG-style embedding and chunking
	err = openai.LoadAndEmbedResumeChunks(resumeData)
	if err != nil {
		log.Fatalf("Failed to embed resume chunks: %v", err)
	}
	fmt.Println("✅ Resume chunks embedded and ready for retrieval")
}


func main() {

	port := config.GetServerPort()
	log.Printf("✅ Server started on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, RegisterRoutes()))
}






