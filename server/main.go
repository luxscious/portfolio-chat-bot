package main

import (
	"fmt"
	"go-ai/config"
	"go-ai/db"
	"go-ai/openai"
	"go-ai/resume"
	"log"
	"net/http"
)

var globalResumeData *resume.ResumeData

func init() {
	config.LoadEnv()
	db.InitMongo()

	resumeData, err := resume.LoadResume("./resume/resume.json")
	if err != nil {
		log.Fatalf("❌ Failed to load resume: %v", err)
	}
	globalResumeData = resumeData

	systemPrompt, err := openai.BuildSystemPrompt(resumeData.PersonaContext)
	if err != nil {
		log.Fatalf("❌ Failed to build system prompt: %v", err)
	}
	fmt.Println("✅ System prompt initialized")
	fmt.Println(systemPrompt)

	fmt.Printf("🔎 %d education entries, %d work entries, %d projects loaded\n",
		len(resumeData.Education), len(resumeData.WorkExperience), len(resumeData.Projects))

	// Getting rid of embedding for now. Using chat gpt to find relevant projects
	// if err := openai.LoadAndEmbedResumeChunks(resumeData); err != nil {
	// 	log.Fatalf("❌ Failed to embed resume chunks: %v", err)
	// }
	// fmt.Println("✅ Resume chunks embedded and ready for retrieval")
}

func main() {
	port := config.GetServerPort()
	log.Printf("✅ Server started on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, RegisterRoutes()))
}
