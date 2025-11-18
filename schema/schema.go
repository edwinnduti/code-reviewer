package schema

import (
	"github.com/tmc/langchaingo/llms"
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
	Score    int     `json:"score"`
}

type CodeReviewer struct {
	Llm      llms.Model
	Template *prompts.PromptTemplate
}
