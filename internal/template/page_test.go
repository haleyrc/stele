package template_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/haleyrc/assert"
	"github.com/haleyrc/stele/internal/site"
	"github.com/haleyrc/stele/internal/template"
)

func TestTemplateRenderer_RenderIndex(t *testing.T) {
	s := newTestSite()
	renderer := template.NewTemplateRenderer()

	var buf bytes.Buffer
	err := renderer.RenderIndex(context.Background(), &buf, s)
	assert.OK(t, err).Fatal()

	compareGolden(t, buf.String(), "testdata/golden/page_homepage.html")
}

func TestTemplateRenderer_RenderIndex_WithSeries(t *testing.T) {
	s := newTestSite()

	// Create a test series
	series := newTestSeries("go-basics", "Go Basics")

	// Update posts to include series posts
	s.Posts = site.Posts{
		newTestPostWithSeries("go-basics/deep-dive", "Deep Dive", "2024-02-01", series),
		newTestPostWithSeries("go-basics/intro", "Introduction", "2024-01-01", series),
		newTestPost("standalone", "Standalone Post", "2024-01-15"),
	}

	// Set up the series with its posts
	series.Posts = site.Posts{
		s.Posts[1], // intro (older, so first in series)
		s.Posts[0], // deep-dive (newer, so second in series)
	}

	// Add series to site
	s.Series = site.AllSeries{series}

	renderer := template.NewTemplateRenderer()

	var buf bytes.Buffer
	err := renderer.RenderIndex(context.Background(), &buf, s)
	assert.OK(t, err).Fatal()

	compareGolden(t, buf.String(), "testdata/golden/page_homepage_with_series.html")
}

func TestTemplateRenderer_RenderTagIndex(t *testing.T) {
	s := newTestSite()
	renderer := template.NewTemplateRenderer()

	var buf bytes.Buffer
	err := renderer.RenderTagIndex(context.Background(), &buf, s)
	assert.OK(t, err).Fatal()

	compareGolden(t, buf.String(), "testdata/golden/page_tag_index.html")
}

func TestTemplateRenderer_RenderArchiveIndex(t *testing.T) {
	s := newTestSite()
	renderer := template.NewTemplateRenderer()

	var buf bytes.Buffer
	err := renderer.RenderArchiveIndex(context.Background(), &buf, s)
	assert.OK(t, err).Fatal()

	compareGolden(t, buf.String(), "testdata/golden/page_archive_index.html")
}
