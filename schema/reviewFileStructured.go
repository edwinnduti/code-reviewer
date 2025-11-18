package schema

import (
	"context"
	"encoding/json"
	"fmt"
	"go/parser"
	"go/token"
	"os"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/prompts"
)

func (cr *CodeReviewer) ReviewFileStructured(filename string) (*ReviewResult, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("reading file: %w", err)
	}

	// Parse for line numbers
	fset := token.NewFileSet()
	_, err = parser.ParseFile(fset, filename, content, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("parsing Go file: %w", err)
	}

	template := prompts.NewPromptTemplate(`
Analyze this Go code and return a JSON response with this exact structure:

{
  "filename": "{{.filename}}",
  "issues": [
    {
      "severity": "critical|warning|suggestion",
      "type": "bug|performance|style|security",
      "line": 42,
      "description": "Detailed issue description",
      "suggestion": "How to fix this issue"
    }
  ],
  "score": 85
}

Code to analyze:
'''go
{{.code}}
'''

Focus on real issues. Score: 100 = perfect, 0 = many serious issues.`,
		[]string{"code", "filename"})

	prompt, err := template.Format(map[string]any{
		"code":     string(content),
		"filename": filename,
	})
	if err != nil {
		return nil, fmt.Errorf("formatting prompt: %w", err)
	}

	ctx := context.Background()
	response, err := cr.Llm.GenerateContent(ctx, []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeHuman, prompt),
	}, llms.WithJSONMode())
	if err != nil {
		return nil, fmt.Errorf("generating review: %w", err)
	}

	var result ReviewResult
	if err := json.Unmarshal([]byte(response.Choices[0].Content), &result); err != nil {
		return nil, fmt.Errorf("parsing JSON response: %w", err)
	}

	return &result, nil
}
