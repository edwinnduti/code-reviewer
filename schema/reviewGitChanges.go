package schema

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

func ReviewGitChanges(reviewer *CodeReviewer) error {
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
