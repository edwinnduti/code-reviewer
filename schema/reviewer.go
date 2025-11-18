package schema

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/tmc/langchaingo/llms/googleai"
	"github.com/tmc/langchaingo/prompts"
)

func NewCodeReviewer(ctx context.Context) (*CodeReviewer, error) {
	// Setup Environment
	fmt.Println("Initializing LLM and Agent...")

	// This requires the 'GEMINI_API_KEY' environment variable to be set.
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		log.Fatal("GEMINI_API_KEY environment variable not set. Please set it to your API key.")
		return nil, fmt.Errorf("GEMINI_API_KEY environment variable not set. Please set it to your API key.")
	}

	// Initialize the Gemini LLM
	llm, err := googleai.New(
		ctx,
		googleai.WithAPIKey(apiKey),
	)
	if err != nil {
		return nil, err
	}

	template := prompts.NewPromptTemplate(`
You are an expert Go code reviewer. Analyze this Go code for:

1. **Bugs and Logic Issues**: Potential runtime errors, nil pointer dereferences, race conditions
2. **Performance**: Inefficient algorithms, unnecessary allocations, string concatenation issues
3. **Style**: Go idioms, naming conventions, error handling patterns
4. **Security**: Input validation, sensitive data handling

Code to review:
'''go
{{.code}}
'''

File: {{.filename}}

Provide specific, actionable feedback. For each issue:
- Explain WHY it's a problem
- Show HOW to fix it with code examples
- Rate severity: Critical, Warning, Suggestion

Focus on the most important issues first.`,
		[]string{"code", "filename"})

	return &CodeReviewer{
		Llm:      llm,
		Template: &template,
	}, nil
}
