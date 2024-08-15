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

func Parse(path string) (*Metadata, error) {
	md := goldmark.New(goldmark.WithExtensions(
		emoji.Emoji,
		extension.GFM,
		&frontmatter.Extender{},
		&katex.Extender{},
	))
	ctx := parser.NewContext()

	contents, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("markdown: parse: %w", err)
	}

	if err := md.Convert(contents, io.Discard, parser.WithContext(ctx)); err != nil {
		return nil, fmt.Errorf("markdown: parse: %w", err)
	}

	var meta Metadata
	if err := frontmatter.Get(ctx).Decode(&meta); err != nil {
		return nil, fmt.Errorf("markdown: parse: %w", err)
	}

	return &meta, nil
}

type Metadata struct {
	Description string   `yaml:"description"`
	Tags        []string `yaml:"tags"`
	Title       string   `yaml:"title"`
}
