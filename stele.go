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

	if err := renderIndex(ctx, dstDir, cfg, posts); err != nil {
		return fmt.Errorf("stele: build: %w", err)
	}

	if err := renderPages(ctx, dstDir, cfg, pages, posts); err != nil {
		return fmt.Errorf("stele: build: %w", err)
	}

	if err := renderPosts(ctx, dstDir, cfg, posts); err != nil {
		return fmt.Errorf("stele: build: %w", err)
	}

	if err := renderArchive(ctx, dstDir, cfg, posts); err != nil {
		return fmt.Errorf("stele: build: %w", err)
	}

	if err := renderTags(ctx, dstDir, cfg, posts); err != nil {
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

func renderArchive(ctx context.Context, dir string, cfg *Config, posts Posts) error {
	path := filepath.Join(dir, "archive.html")

	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("render archive: %w", err)
	}
	defer f.Close()

	postsByYear := posts.ByYear()
	copyright := time.Now().Year()
	if first := posts.First(); first != nil {
		copyright = first.Timestamp.Year()
	}

	vm := template.PostIndexProps{
		Layout: template.LayoutProps{
			Author:      cfg.Author,
			Copyright:   strconv.Itoa(copyright),
			Description: cfg.Description,
			Menu:        menuLinksToProps(cfg.Menu),
			Name:        "Archive",
			Title:       cfg.Title,
		},
		Entries: postIndexToProps(postsByYear),
		Prefix:  "/archive/",
	}

	log.Printf("Rendering %s...", path)
	if err := template.PostIndex(vm).Render(ctx, f); err != nil {
		return fmt.Errorf("render archive: %w", err)
	}

	for _, entry := range postsByYear {
		path := filepath.Join(dir, "archive", entry.Key+".html")

		f, err := os.Create(path)
		if err != nil {
			return fmt.Errorf("render archive: %w", err)
		}

		copyright := time.Now().Year()
		if first := posts.First(); first != nil {
			copyright = first.Timestamp.Year()
		}

		vm := template.PostListProps{
			Layout: template.LayoutProps{
				Author:      cfg.Author,
				Copyright:   strconv.Itoa(copyright),
				Description: cfg.Description,
				Menu:        menuLinksToProps(cfg.Menu),
				Name:        fmt.Sprintf("Posts from %s", entry.Key),
				Title:       cfg.Title,
			},
			Heading: fmt.Sprintf("Posts from %s", entry.Key),
			Posts:   postsToProps(entry.Posts),
		}

		log.Printf("Rendering %s...", path)
		if err := template.PostList(vm).Render(ctx, f); err != nil {
			return fmt.Errorf("render archive: %w", err)
		}
	}

	return nil
}

func renderIndex(ctx context.Context, dir string, cfg *Config, posts Posts) error {
	path := filepath.Join(dir, "index.html")

	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("render index: %w", err)
	}
	defer f.Close()

	latestPost, rest := posts.Head()
	recentPosts := rest.MostRecent(10)

	copyright := time.Now().Year()
	if first := posts.First(); first != nil {
		copyright = first.Timestamp.Year()
	}

	vm := template.IndexProps{
		Layout: template.LayoutProps{
			Author:      cfg.Author,
			Copyright:   strconv.Itoa(copyright),
			Description: cfg.Description,
			Menu:        menuLinksToProps(cfg.Menu),
			Name:        "Home Page",
			Title:       cfg.Title,
		},
		RecentPosts: postsToProps(recentPosts),
	}

	if latestPost != nil {
		postProps := postToProps(*latestPost)
		vm.LatestPost = &postProps
	}

	log.Printf("Rendering %s...", path)
	if err := template.Index(vm).Render(ctx, f); err != nil {
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

func renderPages(ctx context.Context, dir string, cfg *Config, pages Pages, posts Posts) error {
	for _, page := range pages {
		path := filepath.Join(dir, page.Slug+".html")

		f, err := os.Create(path)
		if err != nil {
			return fmt.Errorf("render pages: %w", err)
		}
		defer f.Close()

		copyright := time.Now().Year()
		if first := posts.First(); first != nil {
			copyright = first.Timestamp.Year()
		}

		vm := template.PageProps{
			Layout: template.LayoutProps{
				Author:      cfg.Author,
				Copyright:   strconv.Itoa(copyright),
				Description: cfg.Description,
				Menu:        menuLinksToProps(cfg.Menu),
				Name:        page.Slug,
				Title:       cfg.Title,
			},
			Content: page.Content(),
		}

		log.Printf("Rendering %s...", path)
		if err := template.Page(vm).Render(ctx, f); err != nil {
			return fmt.Errorf("render pages: %w", err)
		}
	}

	return nil
}

func renderPosts(ctx context.Context, dir string, cfg *Config, posts Posts) error {
	for _, post := range posts {
		path := filepath.Join(dir, "posts", post.Slug+".html")

		f, err := os.Create(path)
		if err != nil {
			return fmt.Errorf("render posts: %w", err)
		}
		defer f.Close()

		copyright := time.Now().Year()
		if first := posts.First(); first != nil {
			copyright = first.Timestamp.Year()
		}

		vm := template.PostProps{
			Layout: template.LayoutProps{
				Author:      cfg.Author,
				Copyright:   strconv.Itoa(copyright),
				Description: cfg.Description,
				Menu:        menuLinksToProps(cfg.Menu),
				Name:        post.Title,
				Title:       cfg.Title,
			},
			Post: postToProps(post),
		}

		log.Printf("Rendering %s...", path)
		if err := template.Post(vm).Render(ctx, f); err != nil {
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

func renderTags(ctx context.Context, dir string, cfg *Config, posts Posts) error {
	path := filepath.Join(dir, "tags.html")

	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("render tags: %w", err)
	}
	defer f.Close()

	postsByTag := posts.ByTag()
	copyright := time.Now().Year()
	if first := posts.First(); first != nil {
		copyright = first.Timestamp.Year()
	}

	vm := template.PostIndexProps{
		Layout: template.LayoutProps{
			Author:      cfg.Author,
			Copyright:   strconv.Itoa(copyright),
			Description: cfg.Description,
			Menu:        menuLinksToProps(cfg.Menu),
			Name:        "Tags",
			Title:       cfg.Title,
		},
		Entries: postIndexToProps(postsByTag),
		Prefix:  "/tags/",
	}

	log.Printf("Rendering %s...", path)
	if err := template.PostIndex(vm).Render(ctx, f); err != nil {
		return fmt.Errorf("render tags: %w", err)
	}

	for _, entry := range postsByTag {
		path := filepath.Join(dir, "tags", entry.Key+".html")

		f, err := os.Create(path)
		if err != nil {
			return fmt.Errorf("render tags: %w", err)
		}

		copyright := time.Now().Year()
		if first := posts.First(); first != nil {
			copyright = first.Timestamp.Year()
		}

		vm := template.PostListProps{
			Layout: template.LayoutProps{
				Author:      cfg.Author,
				Copyright:   strconv.Itoa(copyright),
				Description: cfg.Description,
				Menu:        menuLinksToProps(cfg.Menu),
				Name:        fmt.Sprintf("Posts tagged %q", entry.Key),
				Title:       cfg.Title,
			},
			Heading: fmt.Sprintf("Posts tagged %q", entry.Key),
			Posts:   postsToProps(entry.Posts),
		}

		log.Printf("Rendering %s...", path)
		if err := template.PostList(vm).Render(ctx, f); err != nil {
			return fmt.Errorf("render tags: %w", err)
		}
	}

	return nil
}
