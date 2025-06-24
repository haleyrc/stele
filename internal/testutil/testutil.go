// Package testutil provides shared testing utilities for the stele package.
package testutil

import (
	"testing"
	"time"

	"github.com/haleyrc/assert"
	"github.com/haleyrc/stele/internal/site"
)

// AssertRenderedOutput verifies that actual rendered output matches expected output.
// This is used for golden file testing where we compare rendered content against
// known good output.
func AssertRenderedOutput(t *testing.T, expected, actual string) {
	assert.Equal(t, "rendered output", expected, actual)
}

// TestSite returns a canonical site configuration for all tests.
// This ensures consistency across all test suites.
func TestSite() *site.Site {
	return &site.Site{
		Config: site.SiteConfig{
			Author:  "Alice Smith",
			BaseURL: "https://alice.dev",
			Categories: []string{
				"technology",
				"programming",
				"web-development",
			},
			Description: "Alice's development blog covering Go, web APIs, and software engineering",
			Title:       "Alice Codes",
		},
		Posts: site.Posts{
			{
				Frontmatter: site.PostFrontmatter{
					Title:       "Getting Started with Go",
					Description: "A comprehensive beginner's guide to the Go programming language",
					Tags:        []string{"go", "programming", "tutorial"},
					Timestamp:   time.Now(),
				},
				Slug: "getting-started-with-go",
			},
			{
				Frontmatter: site.PostFrontmatter{
					Title:       "Building REST APIs with Go",
					Description: "Learn how to build scalable and maintainable REST APIs using Go",
					Tags:        []string{"go", "api", "web", "rest"},
					Timestamp:   time.Now(),
				},
				Slug: "building-rest-apis-go",
			},
			{
				Frontmatter: site.PostFrontmatter{
					Title:       "Advanced Go Patterns",
					Description: "Exploring advanced design patterns and best practices in Go",
					Tags:        []string{"go", "patterns", "advanced", "best-practices"},
					Timestamp:   time.Now(),
				},
				Slug: "advanced-go-patterns",
			},
			{
				Frontmatter: site.PostFrontmatter{
					Title:       "Testing in Go: A Complete Guide",
					Description: "Comprehensive guide to testing Go applications with examples and best practices",
					Tags:        []string{"go", "testing", "quality", "best-practices"},
					Timestamp:   time.Now(),
				},
				Slug: "testing-in-go-complete-guide",
			},
		},
	}
}

// ExpectedManifest is the expected manifest output for TestSite().
const ExpectedManifest = `{
  "background_color": "white",
  "categories": [
    "technology",
    "programming",
    "web-development"
  ],
  "description": "Alice's development blog covering Go, web APIs, and software engineering",
  "display": "fullscreen",
  "icons": [],
  "name": "Alice Codes",
  "start_url": "https://alice.dev"
}`
