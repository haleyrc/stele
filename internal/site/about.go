package site

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/haleyrc/stele/internal/markdown"
)

// About represents the About page content.
type About struct {
	// The rendered HTML content of the about page.
	Content string
}

// LoadAbout loads the about.md file from the site directory and returns the
// parsed About page. If the about.md file does not exist, returns nil with no
// error.
func LoadAbout(dir string) (*About, error) {
	path := filepath.Join(dir, "about.md")

	if _, err := os.Stat(path); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, fmt.Errorf("load about: %s: %w", path, err)
	}

	var content strings.Builder
	if err := markdown.Parse(path, &content); err != nil {
		return nil, fmt.Errorf("load about: %w", err)
	}

	about := &About{
		Content: content.String(),
	}

	return about, nil
}
