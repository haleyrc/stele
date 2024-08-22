package blog

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Author      string     `yaml:"author"`
	BaseURL     string     `yaml:"baseURL"`
	Categories  []string   `yaml:"categories"`
	Description string     `yaml:"description"`
	Menu        []MenuLink `yaml:"menu"`
	Title       string     `yaml:"title"`
}

func NewConfig(path string) (*Config, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("blog: new config: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(bytes, &cfg); err != nil {
		return nil, fmt.Errorf("blog: new config: %w", err)
	}

	return &cfg, nil
}

type MenuLink struct {
	Label string `yaml:"label"`
	Path  string `yaml:"path"`
}
