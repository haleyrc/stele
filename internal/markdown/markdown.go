// Package markdown contains thin wrappers around third-party markdown
// implementations for the v3 module.
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

// Parse reads the file at path and writes the converted markdown content to w.
func Parse(path string, w io.Writer) error {
	contents, err := os.ReadFile(path) // #nosec G304 - User-specified markdown file is intentional
	if err != nil {
		return fmt.Errorf("markdown: parse: %s: %w", path, err)
	}

	if err := defaultParser.Convert(contents, w); err != nil {
		return fmt.Errorf("markdown: parse: %s: %w", path, err)
	}

	return nil
}

// ParseFrontmatter reads the file at path and attempts to populate fm with the
// values found in the markdown frontmatter block.
func ParseFrontmatter(path string, fm any) error {
	ctx := parser.NewContext()

	contents, err := os.ReadFile(path) // #nosec G304 - User-specified markdown file is intentional
	if err != nil {
		return fmt.Errorf("markdown: parse frontmatter: %s: %w", path, err)
	}

	if err := defaultParser.Convert(contents, io.Discard, parser.WithContext(ctx)); err != nil {
		return fmt.Errorf("markdown: parse frontmatter: %s: %w", path, err)
	}

	if err := frontmatter.Get(ctx).Decode(fm); err != nil {
		return fmt.Errorf("markdown: parse frontmatter: %s: %w", path, err)
	}

	return nil
}
