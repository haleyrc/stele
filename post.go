package stele

import (
	"bytes"
	"cmp"
	"fmt"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/haleyrc/stele/internal/markdown"
	"github.com/haleyrc/stele/template"
	"github.com/haleyrc/stele/template/components"
)

type Post struct {
	Description string
	Path        string
	Slug        string
	Tags        []string
	Timestamp   time.Time
	Title       string
}

func postToProps(p Post) components.PostProps {
	return components.PostProps{
		Content:   p.Content(),
		Slug:      p.Slug,
		Tags:      p.Tags,
		Timestamp: p.Timestamp,
		Title:     p.Title,
	}
}

func NewPost(path string) (*Post, error) {
	meta, err := markdown.ParseFrontmatter(path)
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

func postIndexToProps(index PostIndex) []template.PostIndexEntryProps {
	props := make([]template.PostIndexEntryProps, 0, len(index))
	for _, entry := range index {
		props = append(props, template.PostIndexEntryProps{
			Count: len(entry.Posts),
			Key:   entry.Key,
		})
	}
	return props
}

type PostIndexEntry struct {
	Key   string
	Posts Posts
}

type Posts []Post

func postsToProps(posts Posts) components.PostListProps {
	props := make([]components.PostListEntryProps, 0, len(posts))
	for _, post := range posts {
		props = append(props, components.PostListEntryProps{
			Slug:      post.Slug,
			Timestamp: post.Timestamp,
			Title:     post.Title,
		})
	}
	return components.PostListProps{Posts: props}
}

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

func (ps Posts) Head() (*Post, Posts) {
	if len(ps) == 0 {
		return nil, nil
	}
	if len(ps) == 1 {
		return &ps[0], nil
	}
	return &ps[0], ps[1:]
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
	var buff bytes.Buffer
	if err := markdown.Parse(p.Path, &buff); err != nil {
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