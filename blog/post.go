package blog

import (
	"bytes"
	"cmp"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"time"

	katex "github.com/FurqanSoftware/goldmark-katex"
	"github.com/haleyrc/stele/internal/markdown"
	"github.com/yuin/goldmark"
	emoji "github.com/yuin/goldmark-emoji"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"go.abhg.dev/goldmark/frontmatter"
)

type Post struct {
	Description string
	Path        string
	Slug        string
	Tags        []string
	Timestamp   time.Time
	Title       string
}

func NewPost(path string) (*Post, error) {
	meta, err := markdown.Parse(path)
	if err != nil {
		return nil, fmt.Errorf("blog: new post: %w", err)
	}

	slug, timestamp, err := parsePostPath(path)
	if err != nil {
		return nil, fmt.Errorf("blog: load posts: %w", err)
	}

	post := &Post{
		Description: meta.Description,
		Path:        path,
		Slug:        slug,
		Tags:        meta.Tags,
		Timestamp:   timestamp,
		Title:       meta.Title,
	}

	return post, nil
}

type PostIndex []PostIndexEntry

type PostIndexEntry struct {
	Key   string
	Posts Posts
}

type Posts []Post

func NewPosts(dir string) (Posts, error) {
	files, err := filepath.Glob(filepath.Join(dir, "*.md"))
	if err != nil {
		return nil, fmt.Errorf("blog: new posts: %w", err)
	}

	posts := make(Posts, 0, len(files))
	for _, file := range files {
		post, err := NewPost(file)
		if err != nil {
			return nil, fmt.Errorf("blog: new posts: %w", err)
		}
		posts = append(posts, *post)
	}

	slices.SortFunc(posts, func(a, b Post) int {
		if a.Timestamp == b.Timestamp {
			return cmp.Compare(a.Slug, b.Slug)
		}
		if a.Timestamp.Before(b.Timestamp) {
			return 1
		}
		return -1
	})

	return posts, nil
}

func (ps Posts) First() *Post {
	if len(ps) == 0 {
		return nil
	}
	return &ps[len(ps)-1]
}

func (ps Posts) Last() *Post {
	if len(ps) == 0 {
		return nil
	}
	return &ps[0]
}

func (ps Posts) MostRecent(n int) Posts {
	if len(ps) < n {
		n = len(ps)
	}
	return ps[:n]
}

func (ps Posts) ByTag() PostIndex {
	m := map[string]Posts{}
	for _, post := range ps {
		for _, tag := range post.Tags {
			m[tag] = append(m[tag], post)
		}
	}

	idx := make(PostIndex, 0, len(m))
	for key, posts := range m {
		idx = append(idx, PostIndexEntry{
			Key:   key,
			Posts: posts,
		})
	}

	slices.SortFunc(idx, func(a, b PostIndexEntry) int {
		return cmp.Compare(a.Key, b.Key)
	})

	return idx
}

func (ps Posts) ByYear() PostIndex {
	m := map[string]Posts{}
	for _, post := range ps {
		year := strconv.Itoa(post.Timestamp.Year())
		m[year] = append(m[year], post)
	}

	idx := make(PostIndex, 0, len(m))
	for key, posts := range m {
		idx = append(idx, PostIndexEntry{
			Key:   key,
			Posts: posts,
		})
	}

	slices.SortFunc(idx, func(a, b PostIndexEntry) int {
		return cmp.Compare(b.Key, a.Key)
	})

	return idx
}

func (p *Post) Content() string {
	md := goldmark.New(goldmark.WithExtensions(
		emoji.Emoji,
		extension.GFM,
		&frontmatter.Extender{},
		&katex.Extender{},
	))

	ctx := parser.NewContext()

	contents, err := os.ReadFile(p.Path)
	if err != nil {
		panic(fmt.Errorf("blog: post: content: %w", err))
	}

	var buff bytes.Buffer
	if err := md.Convert(contents, &buff, parser.WithContext(ctx)); err != nil {
		panic(fmt.Errorf("blog: post: content: %w", err))
	}

	return buff.String()
}

func parsePostPath(filename string) (string, time.Time, error) {
	name := strings.TrimSuffix(filepath.Base(filename), ".md")
	nameParts := strings.SplitN(name, "-", 2)

	ts, err := time.Parse("20060102", nameParts[0])
	if err != nil {
		return "", time.Time{}, fmt.Errorf("parse post name: %w", err)
	}

	return nameParts[1], ts, nil
}
