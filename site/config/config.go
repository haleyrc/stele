package config

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
	Menu        []Link
	Name        string
}

func New(path string) (*Config, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, errf("config: new", err)
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
		Name string `yaml:"name"`
	}
	if err := yaml.Unmarshal(bytes, &contents); err != nil {
		return nil, errf("config: new", err)
	}

	cfg := &Config{
		Author:      contents.Author,
		BaseURL:     contents.BaseURL,
		Categories:  contents.Categories,
		Description: contents.Description,
		Name:        contents.Name,
	}

	for _, link := range contents.Menu {
		cfg.Menu = append(cfg.Menu, Link{
			Label: link.Label,
			Path:  link.Path,
		})
	}

	return cfg, nil
}

type Link struct {
	Label string
	Path  string
}

func errf(op string, err error) error {
	return fmt.Errorf("%s: %w", op, err)
}
