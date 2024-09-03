// Package markdown contains thin wrappers around third-party markdown
// implementations.
package markdown

import (
	"fmt"
	"io"
	"os"

	katex "github.com/FurqanSoftware/goldmark-katex"
	"github.com/yuin/goldmark"
	emoji "github.com/yuin/goldmark-emoji"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"go.abhg.dev/goldmark/frontmatter"
)

var defaultParser = goldmark.New(
	goldmark.WithExtensions(
		emoji.Emoji,
		extension.GFM,
		&frontmatter.Extender{},
		&katex.Extender{},
	),
)

// Parse reads the file at path and writes the converted markdown content to
// w.
func Parse(path string, w io.Writer) error {
	contents, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("markdown: parse: %w", err)
	}

	if err := defaultParser.Convert(contents, w); err != nil {
		return fmt.Errorf("markdown: parse: %w", err)
	}

	return nil
}

// Frontmatter represents all of the supported frontmatter fields for posts.
type Frontmatter struct {
	Description string   `yaml:"description"`
	Tags        []string `yaml:"tags"`
	Title       string   `yaml:"title"`
}

// ParseFrontmatter reads the file at path and returns the parsed frontmatter.
func ParseFrontmatter(path string) (*Frontmatter, error) {
	ctx := parser.NewContext()

	contents, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("markdown: parse frontmatter: %w", err)
	}

	if err := defaultParser.Convert(contents, io.Discard, parser.WithContext(ctx)); err != nil {
		return nil, fmt.Errorf("markdown: parse frontmatter: %w", err)
	}

	var fm Frontmatter
	if err := frontmatter.Get(ctx).Decode(&fm); err != nil {
		return nil, fmt.Errorf("markdown: parse frontmatter: %w", err)
	}

	return &fm, nil
}
