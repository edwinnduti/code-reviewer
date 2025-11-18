package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"go/parser"
	"go/token"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/googleai"
	"github.com/tmc/langchaingo/prompts"
)

type Issue struct {
	Severity    string `json:"severity"`
	Type        string `json:"type"`
	Line        int    `json:"line"`
	Description string `json:"description"`
	Suggestion  string `json:"suggestion"`
}

type ReviewResult struct {
	Filename string  `json:"filename"`
	Issues   []Issue `json:"issues"`
	Score    int     `json:"score"` // 0-100
}

type CodeReviewer struct {
	llm      llms.Model
	template *prompts.PromptTemplate
}

func NewCodeReviewer() (*CodeReviewer, error) {
	// Setup Environment and Context
	ctx := context.Background()
	fmt.Println("Initializing LLM and Agent...")

	// This requires the 'GEMINI_API_KEY' environment variable to be set.
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		log.Fatal("GEMINI_API_KEY environment variable not set. Please set it to your API key.")
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
		llm:      llm,
		template: &template,
	}, nil
}

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

	prompt, err := cr.template.Format(map[string]any{
		"code":     string(content),
		"filename": filename,
	})
	if err != nil {
		return fmt.Errorf("formatting prompt: %w", err)
	}

	ctx := context.Background()
	response, err := cr.llm.GenerateContent(ctx, []llms.MessageContent{
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
	response, err := cr.llm.GenerateContent(ctx, []llms.MessageContent{
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

func main() {
	var (
		file = flag.String("file", "", "Go file to review")
		dir  = flag.String("dir", "", "Directory to review (all .go files)")
		git  = flag.Bool("git", false, "Review files changed in git working directory")
	)
	flag.Parse()

	reviewer, err := NewCodeReviewer()
	if err != nil {
		log.Fatal(err)
	}

	switch {
	case *file != "":
		// if err := reviewer.ReviewFile(*file); err != nil {
		response, err := reviewer.ReviewFileStructured(*file)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("\n=== Review for %s ===\n", response.Filename)
		fmt.Println(strings.Repeat("=", 80))

		for _, issue := range response.Issues {
			fmt.Printf("Severity: %s\n", issue.Severity)
			fmt.Printf("Type: %s\n", issue.Type)
			fmt.Printf("Line: %d\n", issue.Line)
			fmt.Printf("Description: %s\n", issue.Description)
			fmt.Printf("Suggestion: %s\n", issue.Suggestion)
			fmt.Println(strings.Repeat("-", 40))
		}

		fmt.Println(strings.Repeat("=", 80))
	case *dir != "":
		if err := reviewDirectory(reviewer, *dir); err != nil {
			log.Fatal(err)
		}
	case *git:
		if err := reviewGitChanges(reviewer); err != nil {
			log.Fatal(err)
		}
	default:
		fmt.Println("Usage:")
		fmt.Println("  code-reviewer -file=main.go")
		fmt.Println("  code-reviewer -dir=./pkg")
		fmt.Println("  code-reviewer -git")
		os.Exit(1)
	}
}

func reviewDirectory(reviewer *CodeReviewer, dir string) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if strings.HasSuffix(path, ".go") && !strings.Contains(path, "vendor/") {
			return reviewer.ReviewFile(path)
		}
		return nil
	})
}

func reviewGitChanges(reviewer *CodeReviewer) error {
	// This is a simplified version - you'd want to use a proper git library
	cmd := exec.Command("git", "diff", "--name-only", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("getting git changes: %w", err)
	}

	files := strings.SplitSeq(strings.TrimSpace(string(output)), "\n")
	for file := range files {
		if strings.HasSuffix(file, ".go") && file != "" {
			if err := reviewer.ReviewFile(file); err != nil {
				log.Printf("Error reviewing %s: %v", file, err)
			}
		}
	}
	return nil
}
