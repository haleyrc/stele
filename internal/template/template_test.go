package template_test

import (
	"flag"
	"os"
	"testing"
	"time"

	"github.com/haleyrc/assert"
	"github.com/haleyrc/stele/internal/site"
)

var update = flag.Bool("update", false, "update golden files")

// compareGolden compares the actual output with the golden file.
// If -update flag is set, it updates the golden file instead.
func compareGolden(t *testing.T, got, goldenPath string) {
	t.Helper()

	if *update {
		err := os.WriteFile(goldenPath, []byte(got), 0644)
		if err != nil {
			t.Fatalf("failed to update golden file: %v", err)
		}
		return
	}

	want, err := os.ReadFile(goldenPath)
	if err != nil {
		t.Fatalf("failed to read golden file: %v", err)
	}

	assert.Equal(t, "output", string(want), got)
}

// newTestSite creates a site with test data for golden path testing.
func newTestSite() *site.Site {
	return &site.Site{
		Config: site.SiteConfig{
			Title:       "Test Blog",
			Description: "A test blog for verification",
			BaseURL:     "https://test.example.com",
			Author:      "Test Author",
		},
		Posts: site.Posts{
			newTestPost("second-post", "Second Post", "2024-01-15"),
			newTestPost("first-post", "First Post", "2024-01-01"),
		},
	}
}

// newTestPost creates a post with test data.
func newTestPost(slug, title, date string) *site.Post {
	return &site.Post{
		Slug: slug,
		Frontmatter: site.PostFrontmatter{
			Title:       title,
			Description: "Test description for " + title,
			Timestamp:   parseDate(date),
			Tags:        []string{"test", "example"},
		},
		Content: "<p>Test content for " + title + "</p>",
	}
}

// newTestPostWithSeries creates a post that's part of a series.
func newTestPostWithSeries(slug, title, date string, series *site.Series) *site.Post {
	post := newTestPost(slug, title, date)
	post.Series = series
	return post
}

// newTestSeries creates a series with the given slug and name.
func newTestSeries(slug, name string) *site.Series {
	return &site.Series{
		Slug: slug,
		Metadata: site.SeriesMetadata{
			Name:        name,
			Description: "Test series description",
		},
	}
}

// parseDate parses a date string in YYYY-MM-DD format.
func parseDate(date string) time.Time {
	t, err := time.Parse("2006-01-02", date)
	if err != nil {
		panic("invalid test date: " + date)
	}
	return t
}
