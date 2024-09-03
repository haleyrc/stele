package stele

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Page represents a page with static HTML content.
type Page struct {
	// The path to the page on disk.
	Path string

	// A URL-safe identifier for the page.
	Slug string
}

// NewPage creates a Page object representing the page at the given path.
func NewPage(path string) (*Page, error) {
	page := &Page{
		Path: path,
		Slug: strings.TrimSuffix(filepath.Base(path), ".html"),
	}
	return page, nil
}

// Content returns the content of the page. This is loaded lazily and this
// method will panic if the file is unavailable or can't be read.
func (p Page) Content() string {
	contents, err := os.ReadFile(p.Path)
	if err != nil {
		panic(fmt.Errorf("blog: page: content: %w", err))
	}
	return string(contents)
}

// Pages is an alias for a slice of Page objects.
type Pages []Page

// NewPages returns a slice of Pages by parsing the contents of the provided
// directory.
func NewPages(dir string) (Pages, error) {
	files, err := filepath.Glob(filepath.Join(dir, "*.html"))
	if err != nil {
		return nil, fmt.Errorf("blog: new pages: %w", err)
	}

	pages := make(Pages, 0, len(files))
	for _, file := range files {
		page, err := NewPage(file)
		if err != nil {
			return nil, fmt.Errorf("blog: new pages: %w", err)
		}
		pages = append(pages, *page)
	}

	return pages, nil
}
