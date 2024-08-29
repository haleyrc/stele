package stele

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Page struct {
	Path string
	Slug string
}

func NewPage(path string) (*Page, error) {
	page := &Page{
		Path: path,
		Slug: strings.TrimSuffix(filepath.Base(path), ".html"),
	}
	return page, nil
}

type Pages []Page

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

func (p Page) Content() string {
	contents, err := os.ReadFile(p.Path)
	if err != nil {
		panic(fmt.Errorf("blog: page: content: %w", err))
	}
	return string(contents)
}
