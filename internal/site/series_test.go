package site_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/haleyrc/assert"
	"github.com/haleyrc/stele/internal/site"
)

func TestLoadSeries(t *testing.T) {
	series, err := site.LoadSeries("testdata/posts/go-basics", false)
	assert.OK(t, err).Fatal()

	assert.Equal(t, "series slug", "go-basics", series.Slug)
	assert.Equal(t, "series name", "Go Basics", series.Metadata.Name)
	assert.Equal(t, "series description", "Learn the basics of Go programming", series.Metadata.Description)
	assert.Equal(t, "post count", 2, len(series.Posts))

	// Posts should be ordered chronologically (oldest first)
	assert.Equal(t, "first post slug", "go-basics/variables", series.Posts[0].Slug)
	assert.Equal(t, "second post slug", "go-basics/functions", series.Posts[1].Slug)
}

func TestLoadSeries_MissingName(t *testing.T) {
	dir := t.TempDir()
	seriesDir := filepath.Join(dir, "test-series")
	if err := os.Mkdir(seriesDir, 0750); err != nil {
		t.Fatal(err)
	}

	// Create index.yaml without name
	indexYAML := `description: "Test description"`
	if err := os.WriteFile(filepath.Join(seriesDir, "index.yaml"), []byte(indexYAML), 0600); err != nil {
		t.Fatal(err)
	}

	_, err := site.LoadSeries(seriesDir, false)
	if err == nil {
		t.Fatal("expected error for series without name")
	}
}

func TestLoadAllSeries(t *testing.T) {
	allSeries, err := site.LoadAllSeries("testdata/posts", false)
	assert.OK(t, err).Fatal()

	// Should find the go-basics series
	if len(allSeries) < 1 {
		t.Fatal("expected at least one series")
	}

	// Verify we can find the go-basics series
	found := false
	for _, s := range allSeries {
		if s.Slug == "go-basics" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected to find go-basics series")
	}
}

func TestAllSeries_GetBySlug(t *testing.T) {
	series1 := &site.Series{Slug: "series-one"}
	series2 := &site.Series{Slug: "series-two"}
	allSeries := site.AllSeries{series1, series2}

	result := allSeries.GetBySlug("series-one")
	assert.Equal(t, "found series", series1, result)

	result = allSeries.GetBySlug("nonexistent")
	if result != nil {
		t.Error("expected nil for nonexistent series")
	}
}

func TestAllSeries_GetSeriesInfo(t *testing.T) {
	post1 := &site.Post{Slug: "test-series/first"}
	post2 := &site.Post{Slug: "test-series/second"}
	post3 := &site.Post{Slug: "test-series/third"}

	series := &site.Series{
		Slug:  "test-series",
		Posts: site.Posts{post1, post2, post3},
	}
	allSeries := site.AllSeries{series}

	// Test middle post
	info := allSeries.GetSeriesInfo("test-series/second")
	if info == nil {
		t.Fatal("expected series info for middle post")
	}
	assert.Equal(t, "series", series, info.Series)
	assert.Equal(t, "position", 2, info.Position)
	assert.Equal(t, "previous", post1, info.Previous)
	assert.Equal(t, "next", post3, info.Next)

	// Test first post
	info = allSeries.GetSeriesInfo("test-series/first")
	if info == nil {
		t.Fatal("expected series info for first post")
	}
	assert.Equal(t, "position", 1, info.Position)
	if info.Previous != nil {
		t.Error("expected no previous post for first post")
	}
	assert.Equal(t, "next", post2, info.Next)

	// Test last post
	info = allSeries.GetSeriesInfo("test-series/third")
	if info == nil {
		t.Fatal("expected series info for last post")
	}
	assert.Equal(t, "position", 3, info.Position)
	assert.Equal(t, "previous", post2, info.Previous)
	if info.Next != nil {
		t.Error("expected no next post for last post")
	}

	// Test standalone post
	info = allSeries.GetSeriesInfo("standalone-post")
	if info != nil {
		t.Error("expected no info for standalone post")
	}
}

func TestAllSeries_AllPosts(t *testing.T) {
	allSeries, err := site.LoadAllSeries("testdata/posts", false)
	assert.OK(t, err).Fatal()

	posts := allSeries.AllPosts()

	// Should have all posts from all series
	if len(posts) < 2 {
		t.Fatal("expected at least 2 posts from series")
	}

	// Posts should be sorted by timestamp descending (newest first)
	// The go-basics series has posts from 2024-01 and 2024-02
	// So the February post should come first
	assert.Equal(t, "first post is functions", "go-basics/functions", posts[0].Slug)
	assert.Equal(t, "second post is variables", "go-basics/variables", posts[1].Slug)
}
