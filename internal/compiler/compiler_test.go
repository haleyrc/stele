package compiler_test

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/haleyrc/assert"
	"github.com/haleyrc/stele/internal/compiler"
	"github.com/haleyrc/stele/internal/site"
)

// mockRenderer tracks all render method calls and writes identifiable content.
type mockRenderer struct {
	mu    sync.Mutex
	calls map[string]int
}

func newMockRenderer() *mockRenderer {
	return &mockRenderer{
		calls: make(map[string]int),
	}
}

func (m *mockRenderer) track(method string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.calls[method]++
}

func (m *mockRenderer) getCalls(method string) int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.calls[method]
}

func (m *mockRenderer) writeContent(w io.Writer, content string) error {
	_, err := fmt.Fprintf(w, "<html><body>%s</body></html>", content)
	return err
}

func (m *mockRenderer) Render404(ctx context.Context, w io.Writer, site *site.Site) error {
	m.track("Render404")
	return m.writeContent(w, "404 Not Found")
}

func (m *mockRenderer) RenderAbout(ctx context.Context, w io.Writer, s *site.Site, about *site.About) error {
	m.track("RenderAbout")
	return m.writeContent(w, "About Page")
}

func (m *mockRenderer) RenderArchiveIndex(ctx context.Context, w io.Writer, site *site.Site) error {
	m.track("RenderArchiveIndex")
	return m.writeContent(w, "Archive Index")
}

func (m *mockRenderer) RenderArchivePage(ctx context.Context, w io.Writer, s *site.Site, year string, posts site.Posts) error {
	m.track("RenderArchivePage")
	return m.writeContent(w, fmt.Sprintf("Archive %s", year))
}

func (m *mockRenderer) RenderIndex(ctx context.Context, w io.Writer, site *site.Site) error {
	m.track("RenderIndex")
	return m.writeContent(w, "Index Page")
}

func (m *mockRenderer) RenderManifest(ctx context.Context, w io.Writer, s *site.Site, manifest *site.Manifest) error {
	m.track("RenderManifest")
	_, err := w.Write([]byte(`{"name":"test"}`))
	return err
}

func (m *mockRenderer) RenderNote(ctx context.Context, w io.Writer, s *site.Site, note *site.Note) error {
	m.track("RenderNote")
	return m.writeContent(w, fmt.Sprintf("Note: %s", note.Slug))
}

func (m *mockRenderer) RenderNotesIndex(ctx context.Context, w io.Writer, site *site.Site) error {
	m.track("RenderNotesIndex")
	return m.writeContent(w, "Notes Index")
}

func (m *mockRenderer) RenderNoteTagIndex(ctx context.Context, w io.Writer, site *site.Site) error {
	m.track("RenderNoteTagIndex")
	return m.writeContent(w, "Note Tag Index")
}

func (m *mockRenderer) RenderNoteTagPage(ctx context.Context, w io.Writer, s *site.Site, tag string, notes site.Notes) error {
	m.track("RenderNoteTagPage")
	return m.writeContent(w, fmt.Sprintf("Note Tag: %s", tag))
}

func (m *mockRenderer) RenderPost(ctx context.Context, w io.Writer, s *site.Site, post *site.Post) error {
	m.track("RenderPost")
	return m.writeContent(w, fmt.Sprintf("Post: %s", post.Slug))
}

func (m *mockRenderer) RenderRSSFeed(ctx context.Context, w io.Writer, s *site.Site, feed *site.RSSFeed) error {
	m.track("RenderRSSFeed")
	_, err := w.Write([]byte(`<?xml version="1.0"?><rss version="2.0"></rss>`))
	return err
}

func (m *mockRenderer) RenderSeriesIndex(ctx context.Context, w io.Writer, s *site.Site, series *site.Series) error {
	m.track("RenderSeriesIndex")
	return m.writeContent(w, fmt.Sprintf("Series: %s", series.Slug))
}

func (m *mockRenderer) RenderTagIndex(ctx context.Context, w io.Writer, site *site.Site) error {
	m.track("RenderTagIndex")
	return m.writeContent(w, "Tag Index")
}

func (m *mockRenderer) RenderTagPage(ctx context.Context, w io.Writer, s *site.Site, tag string, posts site.Posts) error {
	m.track("RenderTagPage")
	return m.writeContent(w, fmt.Sprintf("Tag: %s", tag))
}

// Test helpers

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func assertFileExists(t *testing.T, label, path string) {
	t.Helper()
	if !fileExists(path) {
		t.Errorf("%s: expected file to exist at %s", label, path)
	}
}

func assertFileNotExists(t *testing.T, label, path string) {
	t.Helper()
	if fileExists(path) {
		t.Errorf("%s: expected file to not exist at %s", label, path)
	}
}

// Tests

func TestCompiler_CreatesDirectoryStructure(t *testing.T) {
	testSite, err := site.New("../site/testdata", site.SiteOptions{
		IncludeDrafts:   true,
		NotesExperiment: true,
	})
	assert.OK(t, err).Fatal()

	renderer := newMockRenderer()
	c := compiler.NewCompiler(renderer, testSite)

	outputDir := t.TempDir()
	ctx := context.Background()

	err = c.Compile(ctx, outputDir, "../site/testdata")
	assert.OK(t, err).Fatal()

	// Verify all required directories were created
	dirs := []string{
		outputDir,
		filepath.Join(outputDir, "archive"),
		filepath.Join(outputDir, "notes"),
		filepath.Join(outputDir, "notes", "tags"),
		filepath.Join(outputDir, "posts"),
		filepath.Join(outputDir, "tags"),
	}

	for _, dir := range dirs {
		info, err := os.Stat(dir)
		if err != nil {
			t.Fatalf("expected directory %s to exist: %v", dir, err)
		}
		if !info.IsDir() {
			t.Errorf("expected %s to be a directory", dir)
		}
	}
}

func TestCompiler_CompleteBuild(t *testing.T) {
	testSite, err := site.New("../site/testdata", site.SiteOptions{
		IncludeDrafts:   true,
		NotesExperiment: true,
	})
	assert.OK(t, err).Fatal()

	renderer := newMockRenderer()
	c := compiler.NewCompiler(renderer, testSite)

	outputDir := t.TempDir()
	ctx := context.Background()

	err = c.Compile(ctx, outputDir, "../site/testdata")
	assert.OK(t, err).Fatal()

	// Core pages
	assertFileExists(t, "index page", filepath.Join(outputDir, "index.html"))
	assertFileExists(t, "about page", filepath.Join(outputDir, "about.html"))
	assertFileExists(t, "notes index", filepath.Join(outputDir, "notes.html"))
	assertFileExists(t, "archive index", filepath.Join(outputDir, "archive.html"))
	assertFileExists(t, "tags index", filepath.Join(outputDir, "tags.html"))
	assertFileExists(t, "manifest", filepath.Join(outputDir, "manifest.webmanifest"))
	assertFileExists(t, "rss feed", filepath.Join(outputDir, "rss.xml"))

	// Note pages
	assertFileExists(t, "note tag index", filepath.Join(outputDir, "notes", "tags.html"))
	assertFileExists(t, "golang-tips note", filepath.Join(outputDir, "notes", "golang-tips.html"))
	assertFileExists(t, "algorithms note", filepath.Join(outputDir, "notes", "algorithms.html"))
	assertFileExists(t, "vim-shortcuts note", filepath.Join(outputDir, "notes", "vim-shortcuts.html"))

	// Post pages (including series posts)
	assertFileExists(t, "standalone post", filepath.Join(outputDir, "posts", "getting-started-with-go.html"))
	assertFileExists(t, "draft post", filepath.Join(outputDir, "posts", "draft-exploring-go-generics.html"))

	// Series posts in subdirectories
	assertFileExists(t, "series post 1", filepath.Join(outputDir, "posts", "go-basics", "variables.html"))
	assertFileExists(t, "series post 2", filepath.Join(outputDir, "posts", "go-basics", "functions.html"))

	// Series index page
	assertFileExists(t, "series index", filepath.Join(outputDir, "go-basics.html"))
}

func TestCompiler_ProductionBuildExcludesDrafts(t *testing.T) {
	// Load site WITHOUT drafts (production mode)
	testSite, err := site.New("../site/testdata", site.SiteOptions{IncludeDrafts: false})
	assert.OK(t, err).Fatal()

	renderer := newMockRenderer()
	c := compiler.NewCompiler(renderer, testSite)

	outputDir := t.TempDir()
	ctx := context.Background()

	err = c.Compile(ctx, outputDir, "../site/testdata")
	assert.OK(t, err).Fatal()

	// Draft post should NOT be in the output
	assertFileNotExists(t, "draft post excluded", filepath.Join(outputDir, "posts", "draft-exploring-go-generics.html"))

	// Non-draft posts should still be present
	assertFileExists(t, "non-draft post included", filepath.Join(outputDir, "posts", "getting-started-with-go.html"))

	// Verify the draft post wasn't rendered
	draftCount := 0
	for _, post := range testSite.Posts {
		if strings.Contains(post.Slug, "draft") {
			draftCount++
		}
	}
	assert.Equal(t, "no drafts in production site", 0, draftCount)
}

func TestCompiler_SeriesPostsCreateSubdirectories(t *testing.T) {
	testSite, err := site.New("../site/testdata", site.SiteOptions{IncludeDrafts: true})
	assert.OK(t, err).Fatal()

	renderer := newMockRenderer()
	c := compiler.NewCompiler(renderer, testSite)

	outputDir := t.TempDir()
	ctx := context.Background()

	err = c.Compile(ctx, outputDir, "../site/testdata")
	assert.OK(t, err).Fatal()

	// Verify series subdirectory was created
	seriesDir := filepath.Join(outputDir, "posts", "go-basics")
	info, err := os.Stat(seriesDir)
	if err != nil {
		t.Fatalf("expected series subdirectory to exist: %v", err)
	}
	if !info.IsDir() {
		t.Error("expected go-basics to be a directory")
	}

	// Verify series posts are in subdirectory
	assertFileExists(t, "series post in subdirectory", filepath.Join(seriesDir, "variables.html"))
	assertFileExists(t, "series post in subdirectory", filepath.Join(seriesDir, "functions.html"))

	// Verify standalone posts are NOT in subdirectories (flat structure)
	assertFileExists(t, "standalone post flat", filepath.Join(outputDir, "posts", "getting-started-with-go.html"))
}

func TestCompiler_OptionalAboutPage(t *testing.T) {
	t.Run("with about page", func(t *testing.T) {
		testSite, err := site.New("../site/testdata", site.SiteOptions{IncludeDrafts: false})
		assert.OK(t, err).Fatal()

		renderer := newMockRenderer()
		c := compiler.NewCompiler(renderer, testSite)

		outputDir := t.TempDir()
		ctx := context.Background()

		err = c.Compile(ctx, outputDir, "../site/testdata")
		assert.OK(t, err).Fatal()

		// About page should exist
		assertFileExists(t, "about page created", filepath.Join(outputDir, "about.html"))

		// Renderer should have been called
		assert.Equal(t, "RenderAbout called", 1, renderer.getCalls("RenderAbout"))
	})

	t.Run("without about page", func(t *testing.T) {
		// Create a minimal site without an about page
		testSite := &site.Site{
			Config: site.SiteConfig{
				Title:       "Test Site",
				Description: "Test Description",
				Author:      "Test Author",
				BaseURL:     "https://test.example.com",
			},
			Posts: site.Posts{},
			Notes: site.Notes{},
			About: nil, // No about page
		}

		renderer := newMockRenderer()
		c := compiler.NewCompiler(renderer, testSite)

		outputDir := t.TempDir()
		ctx := context.Background()

		err := c.Compile(ctx, outputDir, ".")
		assert.OK(t, err).Fatal()

		// About page should NOT exist
		assertFileNotExists(t, "about page not created", filepath.Join(outputDir, "about.html"))

		// RenderAbout should NOT have been called
		assert.Equal(t, "RenderAbout not called", 0, renderer.getCalls("RenderAbout"))
	})
}

func TestCompiler_RendererCalledForAllPages(t *testing.T) {
	testSite, err := site.New("../site/testdata", site.SiteOptions{
		IncludeDrafts:   true,
		NotesExperiment: true,
	})
	assert.OK(t, err).Fatal()

	renderer := newMockRenderer()
	c := compiler.NewCompiler(renderer, testSite)

	outputDir := t.TempDir()
	ctx := context.Background()

	err = c.Compile(ctx, outputDir, "../site/testdata")
	assert.OK(t, err).Fatal()

	// Verify core page renders
	assert.Equal(t, "RenderIndex called once", 1, renderer.getCalls("RenderIndex"))
	assert.Equal(t, "RenderAbout called once", 1, renderer.getCalls("RenderAbout"))
	assert.Equal(t, "RenderNotesIndex called once", 1, renderer.getCalls("RenderNotesIndex"))
	assert.Equal(t, "RenderArchiveIndex called once", 1, renderer.getCalls("RenderArchiveIndex"))
	assert.Equal(t, "RenderTagIndex called once", 1, renderer.getCalls("RenderTagIndex"))
	assert.Equal(t, "RenderNoteTagIndex called once", 1, renderer.getCalls("RenderNoteTagIndex"))
	assert.Equal(t, "RenderManifest called once", 1, renderer.getCalls("RenderManifest"))
	assert.Equal(t, "RenderRSSFeed called once", 1, renderer.getCalls("RenderRSSFeed"))

	// Verify note pages rendered (3 notes in testdata)
	assert.Equal(t, "RenderNote called for each note", 3, renderer.getCalls("RenderNote"))

	// Verify post pages rendered (7 posts including series and draft)
	assert.Equal(t, "RenderPost called for each post", 7, renderer.getCalls("RenderPost"))

	// Verify series index rendered (1 series: go-basics)
	assert.Equal(t, "RenderSeriesIndex called for each series", 1, renderer.getCalls("RenderSeriesIndex"))

	// Verify archive pages rendered (depends on years in testdata)
	// testdata has posts from 2024 and 2025
	assert.Equal(t, "RenderArchivePage called for each year", 2, renderer.getCalls("RenderArchivePage"))

	// Verify tag pages rendered (depends on unique tags in testdata)
	postTags := make(map[string]bool)
	for _, post := range testSite.Posts {
		for _, tag := range post.Frontmatter.Tags {
			postTags[tag] = true
		}
	}
	assert.Equal(t, "RenderTagPage called for each unique tag", len(postTags), renderer.getCalls("RenderTagPage"))

	// Verify note tag pages rendered
	noteTags := make(map[string]bool)
	for _, note := range testSite.Notes {
		for _, tag := range note.Frontmatter.Tags {
			noteTags[tag] = true
		}
	}
	assert.Equal(t, "RenderNoteTagPage called for each unique note tag", len(noteTags), renderer.getCalls("RenderNoteTagPage"))
}

func TestCompiler_CleansUpExistingOutputDirectory(t *testing.T) {
	testSite, err := site.New("../site/testdata", site.SiteOptions{IncludeDrafts: false})
	assert.OK(t, err).Fatal()

	renderer := newMockRenderer()
	c := compiler.NewCompiler(renderer, testSite)

	outputDir := t.TempDir()
	ctx := context.Background()

	// Create some existing files in the output directory
	oldFile := filepath.Join(outputDir, "old-file.html")
	err = os.WriteFile(oldFile, []byte("old content"), 0644)
	assert.OK(t, err).Fatal()

	oldDir := filepath.Join(outputDir, "old-directory")
	err = os.Mkdir(oldDir, 0750)
	assert.OK(t, err).Fatal()

	oldNestedFile := filepath.Join(oldDir, "nested-old-file.html")
	err = os.WriteFile(oldNestedFile, []byte("old nested content"), 0644)
	assert.OK(t, err).Fatal()

	// Verify old files exist before compilation
	assertFileExists(t, "old file exists before compile", oldFile)
	assertFileExists(t, "old nested file exists before compile", oldNestedFile)

	// Compile the site
	err = c.Compile(ctx, outputDir, "../site/testdata")
	assert.OK(t, err).Fatal()

	// Verify old files were removed
	assertFileNotExists(t, "old file removed", oldFile)
	assertFileNotExists(t, "old directory removed", oldDir)
	assertFileNotExists(t, "old nested file removed", oldNestedFile)

	// Verify new files were created
	assertFileExists(t, "new index created", filepath.Join(outputDir, "index.html"))
}
