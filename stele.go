package stele

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/haleyrc/stele/template"
)

func Build(ctx context.Context, srcDir, dstDir string) error {
	start := time.Now()

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
	layout := template.Layout{
		Author:      cfg.Author,
		Copyright:   strconv.Itoa(copyright),
		Description: cfg.Description,
		Menu:        menuLinksToProps(cfg.Menu),
		Title:       cfg.Title,
	}

	if err := renderIndex(ctx, dstDir, layout, posts); err != nil {
		return fmt.Errorf("stele: build: %w", err)
	}

	if err := renderPages(ctx, dstDir, layout, pages); err != nil {
		return fmt.Errorf("stele: build: %w", err)
	}

	if err := renderPosts(ctx, dstDir, layout, posts); err != nil {
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

func renderArchive(ctx context.Context, dir string, layout template.Layout, posts Posts) error {
	path := filepath.Join(dir, "archive.html")

	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("render archive: %w", err)
	}
	defer f.Close()

	postsByYear := posts.ByYear()
	props := template.PostIndexProps{
		PageName: "Archive",
		Entries:  postIndexToProps(postsByYear),
		Prefix:   "/archive/",
	}

	log.Printf("Rendering %s...", path)
	if err := layout.PostIndex(props).Render(ctx, f); err != nil {
		return fmt.Errorf("render archive: %w", err)
	}

	for _, entry := range postsByYear {
		path := filepath.Join(dir, "archive", entry.Key+".html")

		f, err := os.Create(path)
		if err != nil {
			return fmt.Errorf("render archive: %w", err)
		}

		props := template.PostListProps{
			Heading: fmt.Sprintf("Posts from %s", entry.Key),
			Posts:   postsToProps(entry.Posts),
		}

		log.Printf("Rendering %s...", path)
		if err := layout.PostList(props).Render(ctx, f); err != nil {
			return fmt.Errorf("render archive: %w", err)
		}
	}

	return nil
}

func renderIndex(ctx context.Context, dir string, layout template.Layout, posts Posts) error {
	path := filepath.Join(dir, "index.html")

	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("render index: %w", err)
	}
	defer f.Close()

	latestPost, rest := posts.Head()
	recentPosts := rest.MostRecent(10)

	props := template.IndexProps{
		RecentPosts: postsToProps(recentPosts),
	}
	if latestPost != nil {
		postProps := postToProps(*latestPost)
		props.LatestPost = &postProps
	}

	log.Printf("Rendering %s...", path)
	if err := layout.Index(props).Render(ctx, f); err != nil {
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

	m, err := NewManifest(cfg)
	if err != nil {
		return err
	}

	log.Printf("Rendering %s...", path)

	return m.Render(ctx, f)
}

func renderPages(ctx context.Context, dir string, layout template.Layout, pages Pages) error {
	for _, page := range pages {
		path := filepath.Join(dir, page.Slug+".html")

		f, err := os.Create(path)
		if err != nil {
			return fmt.Errorf("render pages: %w", err)
		}
		defer f.Close()

		props := template.PageProps{
			Content: page.Content(),
			Slug:    page.Slug,
		}

		log.Printf("Rendering %s...", path)
		if err := layout.Page(props).Render(ctx, f); err != nil {
			return fmt.Errorf("render pages: %w", err)
		}
	}

	return nil
}

func renderPosts(ctx context.Context, dir string, layout template.Layout, posts Posts) error {
	for _, post := range posts {
		path := filepath.Join(dir, "posts", post.Slug+".html")

		f, err := os.Create(path)
		if err != nil {
			return fmt.Errorf("render posts: %w", err)
		}
		defer f.Close()

		props := template.PostProps{
			Post: postToProps(post),
		}

		log.Printf("Rendering %s...", path)
		if err := layout.Post(props).Render(ctx, f); err != nil {
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

	feed, err := NewFeed(cfg, posts)
	if err != nil {
		return err
	}

	log.Printf("Rendering %s...", path)

	return feed.Render(ctx, f)
}

func renderTags(ctx context.Context, dir string, layout template.Layout, posts Posts) error {
	path := filepath.Join(dir, "tags.html")

	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("render tags: %w", err)
	}
	defer f.Close()

	postsByTag := posts.ByTag()
	props := template.PostIndexProps{
		PageName: "Tags",
		Entries:  postIndexToProps(postsByTag),
		Prefix:   "/tags/",
	}

	log.Printf("Rendering %s...", path)
	if err := layout.PostIndex(props).Render(ctx, f); err != nil {
		return fmt.Errorf("render tags: %w", err)
	}

	for _, entry := range postsByTag {
		path := filepath.Join(dir, "tags", entry.Key+".html")

		f, err := os.Create(path)
		if err != nil {
			return fmt.Errorf("render tags: %w", err)
		}

		props := template.PostListProps{
			Heading: fmt.Sprintf("Posts tagged %q", entry.Key),
			Posts:   postsToProps(entry.Posts),
		}

		log.Printf("Rendering %s...", path)
		if err := layout.PostList(props).Render(ctx, f); err != nil {
			return fmt.Errorf("render tags: %w", err)
		}
	}

	return nil
}
