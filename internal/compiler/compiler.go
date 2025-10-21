// Package compiler provides functionality for compiling static site assets.
package compiler

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/haleyrc/stele/internal/site"
)

// Compiler compiles static site assets from a site and renderer.
type Compiler struct {
	// The site to compile.
	Site *site.Site

	// The renderer to use for rendering content.
	Renderer site.Renderer
}

// NewCompiler creates a new compiler for the given site and renderer.
func NewCompiler(renderer site.Renderer, site *site.Site) *Compiler {
	return &Compiler{
		Site:     site,
		Renderer: renderer,
	}
}

// Compile compiles a deployable blog. The resulting assets are written to
// dstDir and source files are read from srcDir. The contents of the destination
// directory, if any, will be deleted when running this function.
func (c *Compiler) Compile(ctx context.Context, dstDir, srcDir string) error {
	if err := c.createOutputDirectory(dstDir); err != nil {
		return fmt.Errorf("build: %w", err)
	}

	if err := c.renderIndexToFile(ctx, dstDir); err != nil {
		return fmt.Errorf("build: %w", err)
	}

	if err := c.renderAboutToFile(ctx, dstDir); err != nil {
		return fmt.Errorf("build: %w", err)
	}

	if len(c.Site.Notes) > 0 {
		if err := c.renderNotesToFiles(ctx, dstDir); err != nil {
			return fmt.Errorf("build: %w", err)
		}

		if err := c.renderNoteTagsToFiles(ctx, dstDir); err != nil {
			return fmt.Errorf("build: %w", err)
		}
	}

	if err := c.renderPostsToFiles(ctx, dstDir); err != nil {
		return fmt.Errorf("build: %w", err)
	}

	if err := c.renderSeriesToFiles(ctx, dstDir); err != nil {
		return fmt.Errorf("build: %w", err)
	}

	if err := c.renderArchiveToFiles(ctx, dstDir); err != nil {
		return fmt.Errorf("build: %w", err)
	}

	if err := c.renderTagsToFiles(ctx, dstDir); err != nil {
		return fmt.Errorf("build: %w", err)
	}

	if err := c.renderManifestToFile(ctx, dstDir); err != nil {
		return fmt.Errorf("build: %w", err)
	}

	if err := c.renderRSSToFile(ctx, dstDir); err != nil {
		return fmt.Errorf("build: %w", err)
	}

	return nil
}

func (c *Compiler) createOutputDirectory(dir string) error {
	if err := os.RemoveAll(dir); err != nil {
		return fmt.Errorf("create output directory: %w", err)
	}

	if err := os.Mkdir(dir, 0750); err != nil {
		return fmt.Errorf("create output directory: %w", err)
	}

	if err := os.Mkdir(filepath.Join(dir, "archive"), 0750); err != nil {
		return fmt.Errorf("create output directory: %w", err)
	}

	if len(c.Site.Notes) > 0 {
		if err := os.Mkdir(filepath.Join(dir, "notes"), 0750); err != nil {
			return fmt.Errorf("create output directory: %w", err)
		}

		if err := os.Mkdir(filepath.Join(dir, "notes", "tags"), 0750); err != nil {
			return fmt.Errorf("create output directory: %w", err)
		}
	}

	if err := os.Mkdir(filepath.Join(dir, "posts"), 0750); err != nil {
		return fmt.Errorf("create output directory: %w", err)
	}

	if err := os.Mkdir(filepath.Join(dir, "tags"), 0750); err != nil {
		return fmt.Errorf("create output directory: %w", err)
	}

	return nil
}

func (c *Compiler) renderToFile(ctx context.Context, path string, renderFn func(context.Context, *os.File) error) error {
	f, err := os.Create(path) // #nosec G304 - User-controlled output directory is intentional
	if err != nil {
		return err
	}
	defer f.Close()

	log.Printf("Writing %s...", path)
	return renderFn(ctx, f)
}

func (c *Compiler) renderIndexToFile(ctx context.Context, dir string) error {
	path := filepath.Join(dir, "index.html")
	return c.renderToFile(ctx, path, func(ctx context.Context, w *os.File) error {
		return c.Renderer.RenderIndex(ctx, w, c.Site)
	})
}

func (c *Compiler) renderAboutToFile(ctx context.Context, dir string) error {
	if c.Site.About == nil {
		return nil
	}

	path := filepath.Join(dir, "about.html")
	return c.renderToFile(ctx, path, func(ctx context.Context, w *os.File) error {
		return c.Renderer.RenderAbout(ctx, w, c.Site, c.Site.About)
	})
}

func (c *Compiler) renderNotesToFiles(ctx context.Context, dir string) error {
	path := filepath.Join(dir, "notes.html")
	if err := c.renderToFile(ctx, path, func(ctx context.Context, w *os.File) error {
		return c.Renderer.RenderNotesIndex(ctx, w, c.Site)
	}); err != nil {
		return fmt.Errorf("render notes: %w", err)
	}

	for _, note := range c.Site.Notes {
		path := filepath.Join(dir, "notes", note.Slug+".html")
		if err := c.renderToFile(ctx, path, func(ctx context.Context, w *os.File) error {
			return c.Renderer.RenderNote(ctx, w, c.Site, note)
		}); err != nil {
			return fmt.Errorf("render notes: %w", err)
		}
	}

	return nil
}

func (c *Compiler) renderNoteTagsToFiles(ctx context.Context, dir string) error {
	path := filepath.Join(dir, "notes", "tags.html")
	if err := c.renderToFile(ctx, path, func(ctx context.Context, w *os.File) error {
		return c.Renderer.RenderNoteTagIndex(ctx, w, c.Site)
	}); err != nil {
		return fmt.Errorf("render note tags: %w", err)
	}

	notesByTag := c.Site.Notes.IndexByTag()
	for _, entry := range notesByTag {
		path := filepath.Join(dir, "notes", "tags", entry.Key+".html")
		if err := c.renderToFile(ctx, path, func(ctx context.Context, w *os.File) error {
			return c.Renderer.RenderNoteTagPage(ctx, w, c.Site, entry.Key, entry.Notes)
		}); err != nil {
			return fmt.Errorf("render note tags: %w", err)
		}
	}

	return nil
}

func (c *Compiler) renderPostsToFiles(ctx context.Context, dir string) error {
	for _, post := range c.Site.Posts {
		// Create subdirectory for series posts if needed
		postDir := filepath.Dir(filepath.Join(dir, "posts", post.Slug+".html"))
		if err := os.MkdirAll(postDir, 0750); err != nil {
			return fmt.Errorf("render posts: %w", err)
		}

		path := filepath.Join(dir, "posts", post.Slug+".html")
		if err := c.renderToFile(ctx, path, func(ctx context.Context, w *os.File) error {
			return c.Renderer.RenderPost(ctx, w, c.Site, post)
		}); err != nil {
			return fmt.Errorf("render posts: %w", err)
		}
	}

	return nil
}

func (c *Compiler) renderSeriesToFiles(ctx context.Context, dir string) error {
	for _, series := range c.Site.Series {
		path := filepath.Join(dir, series.Slug+".html")
		if err := c.renderToFile(ctx, path, func(ctx context.Context, w *os.File) error {
			return c.Renderer.RenderSeriesIndex(ctx, w, c.Site, series)
		}); err != nil {
			return fmt.Errorf("render series: %w", err)
		}
	}

	return nil
}

func (c *Compiler) renderArchiveToFiles(ctx context.Context, dir string) error {
	path := filepath.Join(dir, "archive.html")
	if err := c.renderToFile(ctx, path, func(ctx context.Context, w *os.File) error {
		return c.Renderer.RenderArchiveIndex(ctx, w, c.Site)
	}); err != nil {
		return fmt.Errorf("render archive: %w", err)
	}

	postsByYear := c.Site.Posts.IndexByYear()
	for _, entry := range postsByYear {
		path := filepath.Join(dir, "archive", entry.Key+".html")
		if err := c.renderToFile(ctx, path, func(ctx context.Context, w *os.File) error {
			return c.Renderer.RenderArchivePage(ctx, w, c.Site, entry.Key, entry.Posts)
		}); err != nil {
			return fmt.Errorf("render archive: %w", err)
		}
	}

	return nil
}

func (c *Compiler) renderTagsToFiles(ctx context.Context, dir string) error {
	path := filepath.Join(dir, "tags.html")
	if err := c.renderToFile(ctx, path, func(ctx context.Context, w *os.File) error {
		return c.Renderer.RenderTagIndex(ctx, w, c.Site)
	}); err != nil {
		return fmt.Errorf("render tags: %w", err)
	}

	postsByTag := c.Site.Posts.IndexByTag()
	for _, entry := range postsByTag {
		path := filepath.Join(dir, "tags", entry.Key+".html")
		if err := c.renderToFile(ctx, path, func(ctx context.Context, w *os.File) error {
			return c.Renderer.RenderTagPage(ctx, w, c.Site, entry.Key, entry.Posts)
		}); err != nil {
			return fmt.Errorf("render tags: %w", err)
		}
	}

	return nil
}

func (c *Compiler) renderManifestToFile(ctx context.Context, dir string) error {
	path := filepath.Join(dir, "manifest.webmanifest")
	manifest := c.Site.Manifest()
	return c.renderToFile(ctx, path, func(ctx context.Context, w *os.File) error {
		return c.Renderer.RenderManifest(ctx, w, c.Site, manifest)
	})
}

func (c *Compiler) renderRSSToFile(ctx context.Context, dir string) error {
	path := filepath.Join(dir, "rss.xml")
	feed := c.Site.RSSFeed()
	return c.renderToFile(ctx, path, func(ctx context.Context, w *os.File) error {
		return c.Renderer.RenderRSSFeed(ctx, w, c.Site, feed)
	})
}
