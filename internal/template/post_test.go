package template_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/haleyrc/assert"
	"github.com/haleyrc/stele/internal/site"
	"github.com/haleyrc/stele/internal/template"
)

func TestTemplateRenderer_RenderPost_Standalone(t *testing.T) {
	s := newTestSite()
	post := newTestPost("test-post", "Test Post", "2024-01-01")
	renderer := template.NewTemplateRenderer()

	var buf bytes.Buffer
	err := renderer.RenderPost(context.Background(), &buf, s, post)
	assert.OK(t, err).Fatal()

	compareGolden(t, buf.String(), "testdata/golden/post_standalone.html")
}

func TestTemplateRenderer_RenderPost_SeriesFirst(t *testing.T) {
	s := newTestSite()

	// Create a test series
	series := newTestSeries("tutorial", "Tutorial")

	// Create posts for the series
	part1 := newTestPostWithSeries("tutorial/part-1", "Part 1", "2024-01-01", series)
	part2 := newTestPostWithSeries("tutorial/part-2", "Part 2", "2024-01-15", series)

	// Set up the series with its posts
	series.Posts = site.Posts{part1, part2}

	// Add series to site
	s.Series = site.AllSeries{series}

	renderer := template.NewTemplateRenderer()

	var buf bytes.Buffer
	err := renderer.RenderPost(context.Background(), &buf, s, part1)
	assert.OK(t, err).Fatal()

	compareGolden(t, buf.String(), "testdata/golden/post_series_first.html")
}

func TestTemplateRenderer_RenderPost_SeriesMiddle(t *testing.T) {
	s := newTestSite()

	// Create a test series with three parts
	series := newTestSeries("tutorial", "Tutorial")

	// Create posts for the series
	part1 := newTestPostWithSeries("tutorial/part-1", "Part 1", "2024-01-01", series)
	part2 := newTestPostWithSeries("tutorial/part-2", "Part 2", "2024-01-15", series)
	part3 := newTestPostWithSeries("tutorial/part-3", "Part 3", "2024-02-01", series)

	// Set up the series with its posts
	series.Posts = site.Posts{part1, part2, part3}

	// Add series to site
	s.Series = site.AllSeries{series}

	renderer := template.NewTemplateRenderer()

	var buf bytes.Buffer
	err := renderer.RenderPost(context.Background(), &buf, s, part2)
	assert.OK(t, err).Fatal()

	compareGolden(t, buf.String(), "testdata/golden/post_series_middle.html")
}

func TestTemplateRenderer_RenderPost_SeriesLast(t *testing.T) {
	s := newTestSite()

	// Create a test series
	series := newTestSeries("tutorial", "Tutorial")

	// Create posts for the series
	part1 := newTestPostWithSeries("tutorial/part-1", "Part 1", "2024-01-01", series)
	part2 := newTestPostWithSeries("tutorial/part-2", "Part 2", "2024-01-15", series)

	// Set up the series with its posts
	series.Posts = site.Posts{part1, part2}

	// Add series to site
	s.Series = site.AllSeries{series}

	renderer := template.NewTemplateRenderer()

	var buf bytes.Buffer
	err := renderer.RenderPost(context.Background(), &buf, s, part2)
	assert.OK(t, err).Fatal()

	compareGolden(t, buf.String(), "testdata/golden/post_series_last.html")
}
