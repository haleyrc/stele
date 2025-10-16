package site

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/haleyrc/stele/internal/markdown"
)

// PostFrontmatter represents the YAML frontmatter in a markdown post file.
type PostFrontmatter struct {
	// A short description of the post.
	Description string `yaml:"description"`

	// Whether the post is a draft. Drafts are visible when running the local
	// server, but are not included in production builds.
	Draft bool `yaml:"draft"`

	// A list of tags to associate with the post.
	Tags []string `yaml:"tags"`

	// The "authored date" for the post.
	Timestamp time.Time `yaml:"date"`

	// The title of the post.
	Title string `yaml:"title"`
}

// Validate checks that the frontmatter contains all required fields and that
// field values are valid.
func (fm *PostFrontmatter) Validate() error {
	if fm.Title == "" {
		return fmt.Errorf("posts must have a title")
	}

	if fm.Description == "" {
		return fmt.Errorf("posts must have a description")
	}

	if fm.Draft {
		if !fm.Timestamp.IsZero() {
			return fmt.Errorf("drafts must not have a timestamp")
		}
	} else if fm.Timestamp.IsZero() {
		return fmt.Errorf("posts must have a timestamp")
	}

	return nil
}

// Post represents a blog post.
type Post struct {
	// The YAML frontmatter metadata for the post.
	Frontmatter PostFrontmatter

	// The URL-safe slug for the post. Used to generate post URLs
	// (/posts/{slug}.html).
	Slug string

	// The rendered HTML content of the post.
	Content string

	// The series this post belongs to. Nil for non-series posts.
	Series *Series
}

// LoadPost loads the file at path and returns the parsed post.
func LoadPost(path string) (*Post, error) {
	var fm PostFrontmatter
	if err := markdown.ParseFrontmatter(path, &fm); err != nil {
		return nil, fmt.Errorf("load post: %w", err)
	}

	if err := fm.Validate(); err != nil {
		return nil, fmt.Errorf("load post: %s: %w", path, err)
	}

	if fm.Draft {
		fm.Timestamp = time.Now()
	}

	var content strings.Builder
	if err := markdown.Parse(path, &content); err != nil {
		return nil, fmt.Errorf("load post: %w", err)
	}

	post := &Post{
		Frontmatter: fm,
		Slug:        strings.TrimSuffix(filepath.Base(path), ".md"),
		Content:     content.String(),
	}

	return post, nil
}

// SeriesPosition returns the 1-based position of this post within its series.
// Panics if the post is not part of a series.
func (p *Post) SeriesPosition() int {
	if p.Series == nil {
		panic("SeriesPosition called on non-series post")
	}
	for i, post := range p.Series.Posts {
		if post.Slug == p.Slug {
			return i + 1
		}
	}
	panic(fmt.Sprintf("post %s not found in its series", p.Slug))
}

// Posts is a slice of Post that implements sort.Interface.
// Posts are sorted by timestamp in descending order (newest first).
type Posts []*Post

// LoadPosts loads all markdown files in the given directory and returns the
// parsed posts. If includeDrafts is false, draft posts will be excluded.
func LoadPosts(dir string, includeDrafts bool) (Posts, error) {
	paths, err := filepath.Glob(filepath.Join(dir, "*.md"))
	if err != nil {
		return nil, fmt.Errorf("load posts: %w", err)
	}

	var posts Posts
	for _, path := range paths {
		post, err := LoadPost(path)
		if err != nil {
			return nil, fmt.Errorf("load posts: %w", err)
		}
		if !post.Frontmatter.Draft || includeDrafts {
			posts = append(posts, post)
		}
	}

	posts.Sort()
	return posts, nil
}

// Latest returns the most recent post, or nil if there are no posts.
func (p Posts) Latest() *Post {
	if len(p) == 0 {
		return nil
	}
	return p[0]
}

// Len returns the length of the posts slice.
func (p Posts) Len() int {
	return len(p)
}

// Less reports whether the post at index i should sort before the post at index j.
// Posts are sorted by timestamp in descending order (newest first).
func (p Posts) Less(i, j int) bool {
	return p[i].Frontmatter.Timestamp.After(p[j].Frontmatter.Timestamp)
}

// Recent returns up to maxCount of the most recent posts.
func (p Posts) Recent(maxCount int) Posts {
	if len(p) == 0 {
		return nil
	}

	if maxCount <= 0 || maxCount >= len(p) {
		return p
	}

	return p[:maxCount]
}

// Swap swaps the posts at indices i and j.
func (p Posts) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

// Sort sorts the posts by timestamp in descending order (newest first).
func (p Posts) Sort() {
	sort.Sort(p)
}

// Earliest returns the oldest post, or nil if there are no posts.
func (p Posts) Earliest() *Post {
	if len(p) == 0 {
		return nil
	}
	return p[len(p)-1]
}

// Head returns the first post and the remaining posts.
// If there are no posts, returns nil for both values.
// If there is only one post, returns that post and nil for remaining.
func (p Posts) Head() (*Post, Posts) {
	if len(p) == 0 {
		return nil, nil
	}
	if len(p) == 1 {
		return p[0], nil
	}
	return p[0], p[1:]
}

// IndexByTag returns an index of posts grouped by common tags. A post with
// multiple tags will appear in multiple entries in the index.
func (p Posts) IndexByTag() PostIndex {
	m := map[string]Posts{}
	for _, post := range p {
		for _, tag := range post.Frontmatter.Tags {
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

	sort.Slice(idx, func(i, j int) bool {
		return idx[i].Key < idx[j].Key
	})

	return idx
}

// IndexByYear returns an index of posts grouped by publication year.
func (p Posts) IndexByYear() PostIndex {
	m := map[string]Posts{}
	for _, post := range p {
		year := fmt.Sprintf("%d", post.Frontmatter.Timestamp.Year())
		m[year] = append(m[year], post)
	}

	idx := make(PostIndex, 0, len(m))
	for key, posts := range m {
		idx = append(idx, PostIndexEntry{
			Key:   key,
			Posts: posts,
		})
	}

	sort.Slice(idx, func(i, j int) bool {
		return idx[i].Key > idx[j].Key
	})

	return idx
}

// ForTag returns all posts that have the given tag. If no posts match, returns
// nil.
func (p Posts) ForTag(tag string) Posts {
	var matches Posts
	for _, post := range p {
		for _, t := range post.Frontmatter.Tags {
			if t == tag {
				matches = append(matches, post)
				break
			}
		}
	}
	if len(matches) == 0 {
		return nil
	}
	return matches
}

// ForYear returns all posts from the given year. If no posts match, returns
// nil.
func (p Posts) ForYear(year string) Posts {
	var matches Posts
	for _, post := range p {
		if fmt.Sprintf("%d", post.Frontmatter.Timestamp.Year()) == year {
			matches = append(matches, post)
		}
	}
	if len(matches) == 0 {
		return nil
	}
	return matches
}

// GetBySlug returns the post with the given slug, or nil if not found.
func (p Posts) GetBySlug(slug string) *Post {
	for i := range p {
		if p[i].Slug == slug {
			return p[i]
		}
	}
	return nil
}

// HasTags returns true if any post has non-empty tags.
func (p Posts) HasTags() bool {
	for _, post := range p {
		if len(post.Frontmatter.Tags) > 0 {
			return true
		}
	}
	return false
}

// PostIndex is a slice of entries where each entry contains a set of posts
// that share a common key e.g. post year, tag, etc.
type PostIndex []PostIndexEntry

// PostIndexEntry represents a collection of posts that share a common key e.g.
// post year, tag, etc.
type PostIndexEntry struct {
	// The shared key for the posts (e.g., tag name, year).
	Key string

	// The posts that share this key.
	Posts Posts
}
