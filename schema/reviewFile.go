package schema

import (
	"context"
	"fmt"
	"go/parser"
	"go/token"
	"os"
	"strings"

	"github.com/tmc/langchaingo/llms"
)

func (cr *CodeReviewer) ReviewFile(filename string) error {
	content, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("reading file: %w", err)
	}

	// Parse Go code to ensure it's valid
	fset := token.NewFileSet()
	_, err = parser.ParseFile(fset, filename, content, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("parsing Go file: %w", err)
	}

	prompt, err := cr.Template.Format(map[string]any{
		"code":     string(content),
		"filename": filename,
	})
	if err != nil {
		return fmt.Errorf("formatting prompt: %w", err)
	}

	ctx := context.Background()
	response, err := cr.Llm.GenerateContent(ctx, []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeHuman, prompt),
	})
	if err != nil {
		return fmt.Errorf("generating review: %w", err)
	}

	fmt.Printf("\n=== Review for %s ===\n", filename)
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println(response.Choices[0].Content)
	fmt.Println(strings.Repeat("=", 80))

	return nil
}
