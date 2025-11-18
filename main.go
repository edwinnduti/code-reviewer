package main

import (
	"code-reviewer/schema"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	var (
		file = flag.String("file", "", "Go file to review")
		dir  = flag.String("dir", "", "Directory to review (all .go files)")
		git  = flag.Bool("git", false, "Review files changed in git working directory")
	)
	flag.Parse()

	ctx := context.Background()

	reviewer, err := schema.NewCodeReviewer(ctx)
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
		if err := schema.ReviewDirectory(reviewer, *dir); err != nil {
			log.Fatal(err)
		}
	case *git:
		if err := schema.ReviewGitChanges(reviewer); err != nil {
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
