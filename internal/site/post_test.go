package site_test

import (
	"testing"

	"github.com/haleyrc/assert"
	"github.com/haleyrc/stele/internal/site"
)

func TestLoadPost(t *testing.T) {
	post, err := site.LoadPost("testdata/posts/getting-started-with-go.md")
	assert.OK(t, err).Fatal()

	// Verify frontmatter fields are parsed correctly
	assert.Equal(t, "title", "Getting Started with Go", post.Frontmatter.Title)
	assert.Equal(t, "description", "A comprehensive beginner's guide to the Go programming language", post.Frontmatter.Description)
	assert.Equal(t, "draft", false, post.Frontmatter.Draft)
	assert.Equal(t, "slug", "getting-started-with-go", post.Slug)

	expectedTags := []string{"go", "programming", "tutorial"}
	assert.SliceEqual(t, "tags", expectedTags, post.Frontmatter.Tags)

	// Verify timestamp is parsed (not zero)
	if post.Frontmatter.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp from frontmatter")
	}
}

func TestLoadPost_Draft(t *testing.T) {
	post, err := site.LoadPost("testdata/posts/draft-exploring-go-generics.md")
	assert.OK(t, err).Fatal()

	assert.Equal(t, "title", "Exploring Go Generics", post.Frontmatter.Title)
	assert.Equal(t, "draft", true, post.Frontmatter.Draft)
}

func TestPostsSorting(t *testing.T) {
	// Load posts from testdata
	posts := site.Posts{}

	postFiles := []string{
		"testdata/posts/getting-started-with-go.md",      // 2025-09-20 (newest)
		"testdata/posts/testing-in-go-complete-guide.md", // 2025-09-12 (oldest)
		"testdata/posts/building-rest-apis-go.md",        // 2025-09-18
		"testdata/posts/advanced-go-patterns.md",         // 2025-09-15
	}

	for _, file := range postFiles {
		post, err := site.LoadPost(file)
		assert.OK(t, err).Fatal()
		posts = append(posts, post)
	}

	// Sort posts
	posts.Sort()

	// Verify posts are sorted by timestamp in descending order (newest first)
	assert.Equal(t, "first post slug", "getting-started-with-go", posts[0].Slug)
	assert.Equal(t, "second post slug", "building-rest-apis-go", posts[1].Slug)
	assert.Equal(t, "third post slug", "advanced-go-patterns", posts[2].Slug)
	assert.Equal(t, "fourth post slug", "testing-in-go-complete-guide", posts[3].Slug)

	// Verify timestamps are indeed in descending order
	for i := 0; i < len(posts)-1; i++ {
		if posts[i].Frontmatter.Timestamp.Before(posts[i+1].Frontmatter.Timestamp) {
			t.Errorf("posts[%d] timestamp should be after posts[%d] timestamp", i, i+1)
		}
	}
}

func TestPosts_Head(t *testing.T) {
	// Load and sort posts for testing
	posts := site.Posts{}

	postFiles := []string{
		"testdata/posts/getting-started-with-go.md",      // 2025-09-20 (newest)
		"testdata/posts/building-rest-apis-go.md",        // 2025-09-18
		"testdata/posts/advanced-go-patterns.md",         // 2025-09-15
		"testdata/posts/testing-in-go-complete-guide.md", // 2025-09-12 (oldest)
	}

	for _, file := range postFiles {
		post, err := site.LoadPost(file)
		assert.OK(t, err).Fatal()
		posts = append(posts, post)
	}

	posts.Sort()

	// Test Head with multiple posts
	head, remaining := posts.Head()
	assert.Equal(t, "head post slug", "getting-started-with-go", head.Slug)
	assert.Equal(t, "remaining posts count", 3, len(remaining))
	assert.Equal(t, "first remaining post", "building-rest-apis-go", remaining[0].Slug)
}

func TestPosts_Head_EmptySlice(t *testing.T) {
	posts := site.Posts{}

	head, remaining := posts.Head()
	if head != nil {
		t.Error("expected head to be nil for empty posts")
	}
	if remaining != nil {
		t.Error("expected remaining to be nil for empty posts")
	}
}

func TestPosts_Head_SinglePost(t *testing.T) {
	post, err := site.LoadPost("testdata/posts/getting-started-with-go.md")
	assert.OK(t, err).Fatal()

	posts := site.Posts{post}

	head, remaining := posts.Head()
	assert.Equal(t, "head post slug", "getting-started-with-go", head.Slug)
	if remaining != nil {
		t.Error("expected remaining to be nil for single post")
	}
}

func TestPosts_ByTag(t *testing.T) {
	// Load posts with various tags
	postFiles := []string{
		"testdata/posts/getting-started-with-go.md",      // tags: go, programming, tutorial
		"testdata/posts/building-rest-apis-go.md",        // tags: go, api, web, rest
		"testdata/posts/testing-in-go-complete-guide.md", // tags: go, testing, quality, best-practices
		"testdata/posts/advanced-go-patterns.md",         // tags: go, patterns, advanced, best-practices
	}

	posts := site.Posts{}
	for _, file := range postFiles {
		post, err := site.LoadPost(file)
		assert.OK(t, err).Fatal()
		posts = append(posts, post)
	}

	// Get index by tag
	index := posts.IndexByTag()

	// Verify we have the expected number of unique tags
	// Based on actual testdata: advanced, api, best-practices, go, patterns,
	// programming, quality, rest, testing, tutorial, web
	expectedTagCount := 11
	assert.Equal(t, "tag count", expectedTagCount, len(index))

	// Verify tags are sorted alphabetically
	assert.Equal(t, "first tag", "advanced", index[0].Key)
	assert.Equal(t, "last tag", "web", index[len(index)-1].Key)

	// Find the "go" tag and verify it has all posts
	var goEntry *site.PostIndexEntry
	for i := range index {
		if index[i].Key == "go" {
			goEntry = &index[i]
			break
		}
	}

	if goEntry == nil {
		t.Fatal("expected to find 'go' tag in index")
	}
	assert.Equal(t, "go post count", 4, len(goEntry.Posts))

	// Find the "best-practices" tag and verify it has two posts
	var bpEntry *site.PostIndexEntry
	for i := range index {
		if index[i].Key == "best-practices" {
			bpEntry = &index[i]
			break
		}
	}

	if bpEntry == nil {
		t.Fatal("expected to find 'best-practices' tag in index")
	}
	assert.Equal(t, "best-practices post count", 2, len(bpEntry.Posts))
}

func TestPosts_ByTag_Empty(t *testing.T) {
	posts := site.Posts{}
	index := posts.IndexByTag()

	assert.Equal(t, "empty index", 0, len(index))
}
