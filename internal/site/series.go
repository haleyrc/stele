package site

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

// SeriesMetadata represents the metadata from a series index.yaml file.
type SeriesMetadata struct {
	// The display name for the series.
	Name string `yaml:"name"`

	// An optional description shown on the series index page.
	Description string `yaml:"description"`
}

// Validate checks that the series metadata contains all required fields.
func (sm *SeriesMetadata) Validate() error {
	if sm.Name == "" {
		return fmt.Errorf("series must have a name")
	}
	return nil
}

// Series represents a collection of posts organized together.
type Series struct {
	// The series metadata from index.yaml.
	Metadata SeriesMetadata

	// The URL-safe slug for the series, derived from the directory name.
	Slug string

	// All posts in the series, ordered chronologically (oldest first).
	Posts Posts
}

// LoadSeries loads a series from a directory containing an index.yaml file
// and markdown post files.
func LoadSeries(dir string, includeDrafts bool) (*Series, error) {
	// Load series metadata
	indexPath := filepath.Join(dir, "index.yaml")
	bytes, err := os.ReadFile(indexPath) // #nosec G304 - User-controlled config file path is intentional
	if err != nil {
		return nil, fmt.Errorf("load series: %w", err)
	}

	var metadata SeriesMetadata
	if err := yaml.Unmarshal(bytes, &metadata); err != nil {
		return nil, fmt.Errorf("load series: %w", err)
	}

	if err := metadata.Validate(); err != nil {
		return nil, fmt.Errorf("load series: %w", err)
	}

	// Load all posts in the series directory
	paths, err := filepath.Glob(filepath.Join(dir, "*.md"))
	if err != nil {
		return nil, fmt.Errorf("load series: %w", err)
	}

	var posts Posts
	for _, path := range paths {
		post, err := LoadPost(path)
		if err != nil {
			return nil, fmt.Errorf("load series: %w", err)
		}
		if !post.Frontmatter.Draft || includeDrafts {
			posts = append(posts, post)
		}
	}

	// Sort posts chronologically (oldest first) for series ordering
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].Frontmatter.Timestamp.Before(posts[j].Frontmatter.Timestamp)
	})

	// Update post slugs to include series slug prefix
	slug := filepath.Base(dir)
	for _, post := range posts {
		post.Slug = slug + "/" + post.Slug
	}

	series := &Series{
		Metadata: metadata,
		Slug:     slug,
		Posts:    posts,
	}

	// Set series backlink on each post
	for _, post := range posts {
		post.Series = series
	}

	return series, nil
}

// AllSeries is a slice of Series pointers.
type AllSeries []*Series

// LoadAllSeries discovers and loads all series from the posts directory.
// It looks for subdirectories containing an index.yaml file.
func LoadAllSeries(postsDir string, includeDrafts bool) (AllSeries, error) {
	entries, err := os.ReadDir(postsDir)
	if err != nil {
		return nil, fmt.Errorf("load all series: %w", err)
	}

	var allSeries AllSeries
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		seriesDir := filepath.Join(postsDir, entry.Name())
		indexPath := filepath.Join(seriesDir, "index.yaml")

		// Check if this directory contains a series (has index.yaml)
		if _, err := os.Stat(indexPath); err != nil {
			if os.IsNotExist(err) {
				// Not a series directory, skip
				continue
			}
			return nil, fmt.Errorf("load all series: %w", err)
		}

		series, err := LoadSeries(seriesDir, includeDrafts)
		if err != nil {
			return nil, fmt.Errorf("load all series: %w", err)
		}

		allSeries = append(allSeries, series)
	}

	return allSeries, nil
}

// GetBySlug returns the series with the given slug, or nil if not found.
func (s AllSeries) GetBySlug(slug string) *Series {
	for _, series := range s {
		if series.Slug == slug {
			return series
		}
	}
	return nil
}

// SeriesPostInfo contains information about a post's position within a series.
type SeriesPostInfo struct {
	// The series this post belongs to.
	Series *Series

	// The 1-based position of this post within the series.
	Position int

	// The previous post in the series (nil if this is the first post).
	Previous *Post

	// The next post in the series (nil if this is the last post).
	Next *Post
}

// GetSeriesInfo returns information about where a post fits within its series.
// Returns nil if the post does not belong to a series.
func (s AllSeries) GetSeriesInfo(slug string) *SeriesPostInfo {
	// Check if slug contains a series prefix (e.g., "go-basics/first-post")
	parts := strings.SplitN(slug, "/", 2)
	if len(parts) != 2 {
		return nil // Not a series post
	}

	seriesSlug := parts[0]
	series := s.GetBySlug(seriesSlug)
	if series == nil {
		return nil
	}

	// Find the post within the series
	for i, post := range series.Posts {
		if post.Slug == slug {
			info := &SeriesPostInfo{
				Series:   series,
				Position: i + 1, // 1-based position
			}

			if i > 0 {
				info.Previous = series.Posts[i-1]
			}

			if i < len(series.Posts)-1 {
				info.Next = series.Posts[i+1]
			}

			return info
		}
	}

	return nil
}

// AllPosts returns all posts from all series, sorted by timestamp descending.
func (s AllSeries) AllPosts() Posts {
	var posts Posts
	for _, series := range s {
		posts = append(posts, series.Posts...)
	}
	posts.Sort()
	return posts
}
