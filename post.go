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
	"github.com/haleyrc/stele/template/components"
	"github.com/haleyrc/stele/template/pages"
)

// Frontmatter represents all of the supported frontmatter fields for posts.
type Frontmatter struct {
	// A short description of the post.
	Description string `yaml:"description"`

	// Drafts are visible when running the local server, but are not included in
	// production builds.
	Draft bool `yaml:"draft"`

	// A list of tags to associate with the post.
	Tags []string `yaml:"tags"`

	// The "authored date" for the post.
	Timestamp time.Time `yaml:"date"`

	// The title of the post.
	Title string `yaml:"title"`
}

// Post represents a single markdown post.
type Post struct {
	Frontmatter

	// The path to the file on disk.
	Path string

	// A URL-safe identifier for the post.
	Slug string
}

// NewPost returns a Post object by parsing the file at path.
func NewPost(path string) (*Post, error) {
	var meta Frontmatter
	if err := markdown.ParseFrontmatter(path, &meta); err != nil {
		return nil, fmt.Errorf("blog: new post: %w", err)
	}

	if meta.Title == "" {
		return nil, fmt.Errorf("stele: new post: posts must have a title")
	}

	if meta.Description == "" {
		return nil, fmt.Errorf("stele: new post: posts must have a description")
	}

	if meta.Draft {
		if !meta.Timestamp.IsZero() {
			return nil, fmt.Errorf("stele: new post: drafts can not have a date")
		}
		now := time.Now()
		meta.Timestamp = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	}

	if meta.Timestamp.IsZero() {
		return nil, fmt.Errorf("stele: new post: posts must have a timestamp")
	}

	post := &Post{
		Frontmatter: meta,
		Path:        path,
		Slug:        strings.TrimSuffix(filepath.Base(path), ".md"),
	}

	return post, nil
}

// Content returns the content of the post. This is loaded lazily and this
// method will panic if the file is unavailable or can't be read.
func (p *Post) Content() string {
	var buff bytes.Buffer
	if err := markdown.Parse(p.Path, &buff); err != nil {
		panic(fmt.Errorf("blog: post: content: %w", err))
	}
	return buff.String()
}

// PostIndex is a slice of entries where each entry contains a set of posts that
// share a common key e.g. post year, tag, etc.
type PostIndex []PostIndexEntry

// PostIndexEntry represents a collection of posts that shared a common key e.g.
// post year, tag, etc.
type PostIndexEntry struct {
	Key   string
	Posts Posts
}

// Posts is an alias for a slice of Post objects.
type Posts []Post

// NewPosts returns a slice of Posts by parsing the contents of the provided
// directory.
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

// ByTag returns an index of posts grouped by common tags. A post with multiple
// tags will appear in multiple entries in the index.
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

// ByYear returns an index of posts grouped by the year they were authored. A
// given post will only appear in one entry in the index.
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

// First returns the earliest post by authored date. This method assumes that
// the posts are sorted in descending order by timestamp.
func (ps Posts) First() *Post {
	if len(ps) == 0 {
		return nil
	}
	return &ps[len(ps)-1]
}

// Head returns the first posts in the slice and a Posts object containing the
// remaining posts. If there is only one post, the second return value will be
// nil. If there are no posts, both return values will be nil.
func (ps Posts) Head() (*Post, Posts) {
	if len(ps) == 0 {
		return nil, nil
	}
	if len(ps) == 1 {
		return &ps[0], nil
	}
	return &ps[0], ps[1:]
}

// MostRecent returns the n most recent posts. If n > len(ps), it will return
// all of the posts.
func (ps Posts) MostRecent(n int) Posts {
	if len(ps) < n {
		n = len(ps)
	}
	return ps[:n]
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

func postIndexToProps(index PostIndex) []pages.PostIndexEntryProps {
	props := make([]pages.PostIndexEntryProps, 0, len(index))
	for _, entry := range index {
		props = append(props, pages.PostIndexEntryProps{
			Count: len(entry.Posts),
			Key:   entry.Key,
		})
	}
	return props
}

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
