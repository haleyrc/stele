package stele

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/a-h/templ"

	"github.com/haleyrc/stele/template"
	"github.com/haleyrc/stele/template/pages"
)

type buildConfig struct {
	drafts bool
}

type BuildOpt func(*buildConfig)

func WithDrafts(cfg *buildConfig) {
	cfg.drafts = true
}

// Build compiles a deployable blog. Source files are read from srcDir and the
// resulting assets are written to dstDir. The contents of the destination
// directory, if any, will be deleted when running this function.
func Build(ctx context.Context, srcDir, dstDir string, opts ...BuildOpt) error {
	start := time.Now()

	buildCfg := buildConfig{
		drafts: false,
	}
	for _, opt := range opts {
		opt(&buildCfg)
	}

	log.Println("Loading configuration...")
	cfg, err := NewConfig(filepath.Join(srcDir, "config.yml"))
	if err != nil {
		return fmt.Errorf("stele: build: %w", err)
	}

	log.Println("Indexing posts...")
	posts, err := NewPosts(filepath.Join(srcDir, "posts"))
	if err != nil {
		return fmt.Errorf("stele: build: %w", err)
	}
	log.Printf("Found %d posts.", len(posts))

	log.Println("Indexing pages...")
	pages, err := NewPages(filepath.Join(srcDir, "pages"))
	if err != nil {
		return fmt.Errorf("stele: build: %w", err)
	}
	log.Printf("Found %d pages.", len(pages))

	if err := createOutputDirectory(dstDir); err != nil {
		return fmt.Errorf("stele: build: %w", err)
	}

	copyright := time.Now().Year()
	if first := posts.First(); first != nil {
		copyright = first.Timestamp.Year()
	}
	layout := template.DefaultLayout(
		cfg.Title,
		cfg.Description,
		cfg.Author,
		strconv.Itoa(copyright),
		menuLinksToProps(cfg.Menu),
	)

	if err := renderIndex(ctx, dstDir, layout, posts); err != nil {
		return fmt.Errorf("stele: build: %w", err)
	}

	if err := renderPages(ctx, dstDir, layout, pages); err != nil {
		return fmt.Errorf("stele: build: %w", err)
	}

	if err := renderPosts(ctx, dstDir, layout, posts, buildCfg.drafts); err != nil {
		return fmt.Errorf("stele: build: %w", err)
	}

	if err := renderArchive(ctx, dstDir, layout, posts); err != nil {
		return fmt.Errorf("stele: build: %w", err)
	}

	if err := renderTags(ctx, dstDir, layout, posts); err != nil {
		return fmt.Errorf("stele: build: %w", err)
	}

	if err := renderManifest(ctx, dstDir, cfg); err != nil {
		return fmt.Errorf("stele: build: %w", err)
	}

	if err := renderRSS(ctx, dstDir, cfg, posts); err != nil {
		return fmt.Errorf("stele: build: %w", err)
	}

	log.Printf("Took %s.", time.Since(start))

	return nil
}

func createOutputDirectory(dir string) error {
	if err := os.RemoveAll(dir); err != nil {
		return fmt.Errorf("stele: build: %w", err)
	}

	if err := os.Mkdir(dir, os.ModePerm); err != nil {
		return fmt.Errorf("stele: build: %w", err)
	}

	if err := os.Mkdir(filepath.Join(dir, "archive"), os.ModePerm); err != nil {
		return fmt.Errorf("stele: build: %w", err)
	}

	if err := os.Mkdir(filepath.Join(dir, "posts"), os.ModePerm); err != nil {
		return fmt.Errorf("stele: build: %w", err)
	}

	if err := os.Mkdir(filepath.Join(dir, "tags"), os.ModePerm); err != nil {
		return fmt.Errorf("stele: build: %w", err)
	}

	return nil
}

func renderArchive(ctx context.Context, dir string, layout pages.LayoutFunc, posts Posts) error {
	path := filepath.Join(dir, "archive.html")
	postsByYear := posts.ByYear()

	props := pages.PostIndexProps{
		PageName: "Archive",
		Entries:  postIndexToProps(postsByYear),
		Prefix:   "/archive/",
	}

	if err := renderToPath(ctx, path, pages.PostIndex(layout, props)); err != nil {
		return fmt.Errorf("render archive: %w", err)
	}

	for _, entry := range postsByYear {
		path := filepath.Join(dir, "archive", entry.Key+".html")

		props := pages.PostListProps{
			Heading: fmt.Sprintf("Posts from %s", entry.Key),
			Posts:   postsToProps(entry.Posts),
		}

		if err := renderToPath(ctx, path, pages.PostList(layout, props)); err != nil {
			return fmt.Errorf("render archive: %w", err)
		}
	}

	return nil
}

func renderIndex(ctx context.Context, dir string, layout pages.LayoutFunc, posts Posts) error {
	path := filepath.Join(dir, "index.html")

	latestPost, rest := posts.Head()
	recentPosts := rest.MostRecent(10)

	props := pages.IndexProps{
		RecentPosts: postsToProps(recentPosts),
	}
	if latestPost != nil {
		postProps := postToProps(*latestPost)
		props.LatestPost = &postProps
	}

	if err := renderToPath(ctx, path, pages.Index(layout, props)); err != nil {
		return fmt.Errorf("render index: %w", err)
	}

	return nil
}

func renderManifest(ctx context.Context, dir string, cfg *Config) error {
	path := filepath.Join(dir, "manifest.webmanifest")

	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("render manifest: %w", err)
	}
	defer f.Close()

	log.Printf("Rendering %s...", path)
	return RenderManifest(ctx, f, cfg)
}

func renderPages(ctx context.Context, dir string, layout pages.LayoutFunc, ps []Page) error {
	for _, page := range ps {
		path := filepath.Join(dir, page.Slug+".html")

		props := pages.PageProps{
			Content: page.Content(),
			Slug:    page.Slug,
		}

		if err := renderToPath(ctx, path, pages.Page(layout, props)); err != nil {
			return fmt.Errorf("render pages: %w", err)
		}
	}

	return nil
}

func renderPosts(ctx context.Context, dir string, layout pages.LayoutFunc, posts Posts, renderDrafts bool) error {
	for _, post := range posts {
		if post.Draft && !renderDrafts {
			continue
		}

		path := filepath.Join(dir, "posts", post.Slug+".html")
		props := pages.PostProps{
			Post: postToProps(post),
		}
		if err := renderToPath(ctx, path, pages.Post(layout, props)); err != nil {
			return fmt.Errorf("render posts: %w", err)
		}
	}

	return nil
}

func renderRSS(ctx context.Context, dir string, cfg *Config, posts Posts) error {
	path := filepath.Join(dir, "rss.xml")

	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("render rss: %w", err)
	}
	defer f.Close()

	log.Printf("Rendering %s...", path)
	return RenderRSSFeed(ctx, f, cfg, posts)
}

func renderTags(ctx context.Context, dir string, layout pages.LayoutFunc, posts Posts) error {
	path := filepath.Join(dir, "tags.html")
	postsByTag := posts.ByTag()

	props := pages.PostIndexProps{
		PageName: "Tags",
		Entries:  postIndexToProps(postsByTag),
		Prefix:   "/tags/",
	}

	if err := renderToPath(ctx, path, pages.PostIndex(layout, props)); err != nil {
		return fmt.Errorf("render tags: %w", err)
	}

	for _, entry := range postsByTag {
		path := filepath.Join(dir, "tags", entry.Key+".html")

		props := pages.PostListProps{
			Heading: fmt.Sprintf("Posts tagged %q", entry.Key),
			Posts:   postsToProps(entry.Posts),
		}

		if err := renderToPath(ctx, path, pages.PostList(layout, props)); err != nil {
			return fmt.Errorf("render tags: %w", err)
		}
	}

	return nil
}

func renderToPath(ctx context.Context, path string, component templ.Component) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	log.Printf("Rendering %s...", path)
	return component.Render(ctx, f)
}
