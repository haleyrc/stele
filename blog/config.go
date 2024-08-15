package blog

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Author      string
	BaseURL     string
	Categories  []string
	Description string
	Menu        []MenuLink
	Title       string
}

func NewConfig(path string) (*Config, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("blog: new config: %w", err)
	}

	var contents struct {
		Author      string   `yaml:"author"`
		BaseURL     string   `yaml:"baseURL"`
		Categories  []string `yaml:"categories"`
		Description string   `yaml:"description"`
		Menu        []struct {
			Label string `yaml:"label"`
			Path  string `yaml:"path"`
		} `yaml:"menu"`
		Title string `yaml:"title"`
	}
	if err := yaml.Unmarshal(bytes, &contents); err != nil {
		return nil, fmt.Errorf("blog: new config: %w", err)
	}

	cfg := &Config{
		Author:      contents.Author,
		BaseURL:     contents.BaseURL,
		Categories:  contents.Categories,
		Description: contents.Description,
		Title:       contents.Title,
	}

	for _, link := range contents.Menu {
		cfg.Menu = append(cfg.Menu, MenuLink{
			Label: link.Label,
			Path:  link.Path,
		})
	}

	return cfg, nil
}

type MenuLink struct {
	Label string
	Path  string
}
