package aiagent

import (
	// import standard libraries
	// Import the GenerativeAI package for Go
	"context"
	"log"
	"os"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

var aiagent *genai.Client

// InitializeAIAgent initializes the Generative AI client

func InitializeAIAgent(ctx context.Context) {
	var err error
	aiagent, err = genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GENAI_API_KEY")))
	if err != nil {
		log.Fatal(err)
	}
}

func GetAIAgent() *genai.Client {
	return aiagent
}
