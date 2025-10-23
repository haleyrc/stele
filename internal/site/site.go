package site

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

// SiteOptions contains configuration options for creating a site.
type SiteOptions struct {
	// Whether to include draft posts in the site.
	IncludeDrafts bool

	// Whether to enable the experimental notes feature.
	NotesExperiment bool
}

// Site represents a complete blog site with configuration and content.
type Site struct {
	// The optional About page for the site. Will be nil if no about.md exists.
	About *About

	// The site configuration loaded from stele.yaml.
	Config SiteConfig

	// The root directory of the site content.
	Dir string

	// All notes for the site.
	Notes Notes

	// All posts for the site (includes both standalone posts and series posts).
	Posts Posts

	// All series for the site.
	Series AllSeries

	// The options used when creating the site.
	Opts SiteOptions
}

// New creates and initializes a new Site from the given directory and options.
func New(dir string, opts SiteOptions) (*Site, error) {
	siteStart := time.Now()

	s := &Site{
		Dir:  dir,
		Opts: opts,
	}

	dur, err := logPhase("Loading config", s.loadConfig)
	if err != nil {
		return nil, fmt.Errorf("new site: %w", err)
	}
	log.Printf("Loaded config (%v)", dur)

	dur, err = logPhase("Loading about page", s.loadAbout)
	if err != nil {
		return nil, fmt.Errorf("new site: %w", err)
	}
	log.Printf("Loaded about page (%v)", dur)

	if s.Opts.NotesExperiment {
		dur, err = logPhase("Loading notes", s.loadNotes)
		if err != nil {
			return nil, fmt.Errorf("new site: %w", err)
		}
		log.Printf("Loaded %d notes (%v)", len(s.Notes), dur)
	}

	dur, err = logPhase("Loading series", s.loadSeries)
	if err != nil {
		return nil, fmt.Errorf("new site: %w", err)
	}
	log.Printf("Loaded %d series (%v)", len(s.Series), dur)

	dur, err = logPhase("Loading posts", s.loadPosts)
	if err != nil {
		return nil, fmt.Errorf("new site: %w", err)
	}
	log.Printf("Loaded %d posts (%v)", len(s.Posts), dur)

	log.Printf("Site loaded successfully (%v)", time.Since(siteStart).Round(time.Millisecond))

	return s, nil
}

func logPhase(msg string, fn func() error) (time.Duration, error) {
	log.Printf("%s...", msg)
	start := time.Now()
	err := fn()
	return time.Since(start).Round(time.Millisecond), err
}

func (s *Site) loadConfig() error {
	var err error

	config, err := LoadSiteConfig(s.Dir)
	if err != nil {
		return fmt.Errorf("site: load config: %w", err)
	}
	s.Config = *config

	return nil
}

func (s *Site) loadAbout() error {
	about, err := LoadAbout(s.Dir)
	if err != nil {
		return fmt.Errorf("site: load about: %w", err)
	}
	s.About = about
	return nil
}

func (s *Site) loadNotes() error {
	notes, err := LoadNotes(filepath.Join(s.Dir, "notes"))
	if err != nil {
		return fmt.Errorf("site: load notes: %w", err)
	}
	s.Notes = notes
	return nil
}

func (s *Site) loadSeries() error {
	series, err := LoadAllSeries(filepath.Join(s.Dir, "posts"), s.Opts.IncludeDrafts)
	if err != nil {
		return fmt.Errorf("site: load series: %w", err)
	}
	s.Series = series
	return nil
}

func (s *Site) loadPosts() error {
	// Load standalone posts (markdown files at the root of posts/)
	posts, err := LoadPosts(filepath.Join(s.Dir, "posts"), s.Opts.IncludeDrafts)
	if err != nil {
		return fmt.Errorf("site: load posts: %w", err)
	}

	// Merge series posts with standalone posts
	seriesPosts := s.Series.AllPosts()
	s.Posts = append(posts, seriesPosts...)
	s.Posts.Sort()

	return nil
}

// CopyrightYear returns the year of the earliest post, or the current year if
// no posts exist.
func (s *Site) CopyrightYear() int {
	if first := s.Posts.Earliest(); first != nil {
		return first.Frontmatter.Timestamp.Year()
	}
	return time.Now().Year()
}

// HasSocialLinks returns true if the site has any social media links configured.
func (s *Site) HasSocialLinks() bool {
	return s.Config.Social.GitHub != "" || s.Config.Social.LinkedIn != ""
}

// Manifest creates and returns the web manifest for the site.
func (s *Site) Manifest() *Manifest {
	return NewManifest(s)
}

// RSSFeed creates and returns the RSS feed for the site.
func (s *Site) RSSFeed() *RSSFeed {
	return NewRSSFeed(s)
}

// SiteConfig represents the configuration loaded from stele.yaml.
type SiteConfig struct {
	// The author of the blog.
	Author string `yaml:"author"`

	// The URL where the blog will be hosted. This is used to construct links
	// to resources e.g. in the RSS feed and must be an absolute URL including
	// protocol.
	BaseURL string `yaml:"baseURL"`

	// A list of categories that describe the content of the blog. This is
	// used in the web manifest and RSS feed and should be descriptive without
	// being overloaded.
	Categories []string `yaml:"categories"`

	// A description of the blog's content and/or purpose.
	Description string `yaml:"description"`

	// Social media links to display on the About page.
	Social SocialLinks `yaml:"social"`

	// The title/name of the blog.
	Title string `yaml:"title"`
}

// SocialLinks contains URLs for social media profiles.
type SocialLinks struct {
	// GitHub profile URL.
	GitHub string `yaml:"github"`

	// LinkedIn profile URL.
	LinkedIn string `yaml:"linkedin"`
}

// Validate checks that the site configuration contains all required fields and
// that field values are valid.
func (c *SiteConfig) Validate() error {
	if c.Author == "" {
		return fmt.Errorf("site config must have an author")
	}

	if c.BaseURL == "" {
		return fmt.Errorf("site config must have a base URL")
	}

	u, err := url.Parse(c.BaseURL)
	if err != nil {
		return fmt.Errorf("site config base URL is invalid: %w", err)
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("site config base URL must use http or https scheme")
	}

	if c.Description == "" {
		return fmt.Errorf("site config must have a description")
	}

	if c.Title == "" {
		return fmt.Errorf("site config must have a title")
	}

	return nil
}

// LoadSiteConfig loads the file at path and returns the parsed configuration.
func LoadSiteConfig(dir string) (*SiteConfig, error) {
	path := filepath.Join(dir, "stele.yaml")
	bytes, err := os.ReadFile(path) // #nosec G304 - User-controlled config file path is intentional
	if err != nil {
		return nil, fmt.Errorf("load site config: %s: %w", path, err)
	}

	var cfg SiteConfig
	if err := yaml.Unmarshal(bytes, &cfg); err != nil {
		return nil, fmt.Errorf("load site config: %s: %w", path, err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("load site config: %s: %w", path, err)
	}

	return &cfg, nil
}
