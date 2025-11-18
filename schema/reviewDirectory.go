package schema

import (
	"os"
	"path/filepath"
	"strings"
)

func ReviewDirectory(reviewer *CodeReviewer, dir string) error {
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
