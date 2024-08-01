// Package index provides support for parsing a directory containing blog
// content.
package index

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	katex "github.com/FurqanSoftware/goldmark-katex"
	"github.com/yuin/goldmark"
	emoji "github.com/yuin/goldmark-emoji"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"go.abhg.dev/goldmark/frontmatter"
)

type Index struct {
	Pages []Page
	Posts Posts
}

func New(dir string) (*Index, error) {
	pages, err := indexPages(dir)
	if err != nil {
		return nil, errf("index: new", err)
	}

	posts, err := indexPosts(dir)
	if err != nil {
		return nil, errf("index: new", err)
	}

	return &Index{Pages: pages, Posts: posts}, nil
}

type Page struct {
	Path string
	Slug string
}

func indexPages(dir string) ([]Page, error) {
	files, err := filepath.Glob(filepath.Join(dir, "pages", "*.html"))
	if err != nil {
		return nil, errf("index pages", err)
	}

	log.Println("Indexing pages:")
	pages := make([]Page, 0, len(files))
	for _, file := range files {
		pages = append(pages, Page{
			Path: file,
			Slug: strings.TrimSuffix(filepath.Base(file), ".html"),
		})
		log.Printf("\t%s", file)
	}

	return pages, nil
}

type Post struct {
	Description string
	Path        string
	Slug        string
	Tags        []string
	Timestamp   time.Time
	Title       string
}

type Posts []Post

func indexPosts(dir string) (Posts, error) {
	md := goldmark.New(goldmark.WithExtensions(
		emoji.Emoji,
		extension.GFM,
		&frontmatter.Extender{},
		&katex.Extender{},
	))

	files, err := filepath.Glob(filepath.Join(dir, "posts", "*.md"))
	if err != nil {
		return nil, errf("index posts", err)
	}

	log.Println("Indexing posts:")
	posts := make([]Post, 0, len(files))
	for _, file := range files {
		ctx := parser.NewContext()

		contents, err := os.ReadFile(file)
		if err != nil {
			return nil, errf("load posts", err)
		}

		var buff bytes.Buffer
		if err := md.Convert(contents, &buff, parser.WithContext(ctx)); err != nil {
			return nil, errf("load posts", err)
		}

		var meta struct {
			Description string   `yaml:"description"`
			Tags        []string `yaml:"tags"`
			Title       string   `yaml:"title"`
		}
		if err := frontmatter.Get(ctx).Decode(&meta); err != nil {
			return nil, errf("load posts", err)
		}

		slug, timestamp, err := parsePostName(file)
		if err != nil {
			return nil, errf("load posts", err)
		}

		posts = append(posts, Post{
			Description: meta.Description,
			Path:        file,
			Slug:        slug,
			Tags:        meta.Tags,
			Timestamp:   timestamp,
			Title:       meta.Title,
		})

		log.Printf("\t%s", file)
	}

	return posts, nil
}

// TODO: This isn't very efficient.
// TODO: This will be basically broken for a new site until there is at least one post.
func (ps Posts) First() Post {
	var latest *Post
	for _, p := range ps {
		if latest == nil || p.Timestamp.Before(latest.Timestamp) {
			latest = &p
		}
	}
	return *latest
}

func errf(op string, err error) error {
	return fmt.Errorf("%s: %w", op, err)
}

func parsePostName(filename string) (string, time.Time, error) {
	name := strings.TrimSuffix(filepath.Base(filename), ".md")
	nameParts := strings.SplitN(name, "-", 2)

	ts, err := time.Parse("20060102", nameParts[0])
	if err != nil {
		return "", time.Time{}, errf("parse post name", err)
	}

	return nameParts[1], ts, nil
}
