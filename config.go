package stele

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"

	"github.com/haleyrc/stele/template"
)

// Config represents the supported configuration for the stele framework.
type Config struct {
	// The author of the blog.
	Author string `yaml:"author"`

	// The URL where the blog will be hosted. This is used to construct links
	// to resources e.g. in the RSS feed and must be an absolute URL including
	// protocol.
	BaseURL string `yaml:"baseURL"`

	// A list of categories that describe the content of the blog. This is used in
	// the web manifest and RSS feed and should be descriptive without being
	// overloaded.
	Categories []string `yaml:"categories"`

	// A description of the blog's content and/or purpose.
	Description string `yaml:"description"`

	// A list of links to include in the top-level site navigation. This is useful
	// for including links to raw pages e.g. an "about me" page.
	Menu []MenuLink `yaml:"menu"`

	// The title/name of the blog.
	Title string `yaml:"title"`
}

// NewConfig loads the file at path and returns the parsed configuration.
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

// MenuLink represents a link that will appear in the main page navigation.
type MenuLink struct {
	// The text to show the user.
	Label string `yaml:"label"`

	// The path for the link.
	Path string `yaml:"path"`
}

func menuLinksToProps(links []MenuLink) []template.MenuLink {
	props := make([]template.MenuLink, 0, len(links))
	for _, link := range links {
		props = append(props, template.MenuLink{
			Label: link.Label,
			Path:  link.Path,
		})
	}
	return props
}
