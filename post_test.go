package stele_test

import (
	"path/filepath"
	"slices"
	"testing"

	"github.com/haleyrc/stele"
)

func TestNewPost(t *testing.T) {
	testcases := []struct {
		filename    string
		title       string
		description string
		tags        []string
		draft       bool
	}{
		{
			filename:    "20220103-first-post.md",
			title:       "First Post",
			description: "The first post",
			tags:        []string{"go", "react"},
			draft:       false,
		},
		{
			filename:    "20240406-second-post.md",
			title:       "Second Post",
			description: "The second post",
			tags:        []string{"go", "react"},
			draft:       false,
		},
		{
			filename:    "20240406-third-post.md",
			title:       "Third Post",
			description: "The third post",
			tags:        []string{"go", "react"},
			draft:       true,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.filename, func(t *testing.T) {
			path := filepath.Join("testdata", "posts", tc.filename)
			post, err := stele.NewPost(path)
			if err != nil {
				t.Fatal(err)
			}

			if post.Title != tc.title {
				t.Errorf("expected post.Title = %q, but it was %q", tc.title, post.Title)
			}
			if post.Description != tc.description {
				t.Errorf("expected post.Description = %q, but it was %q", tc.description, post.Description)
			}
			if !slices.Equal(tc.tags, post.Tags) {
				t.Errorf("expected post.Tags = %v, but it was %v", tc.tags, post.Tags)
			}
			if post.Draft != tc.draft {
				t.Errorf("expected post.Draft = %t, but it was %t", tc.draft, post.Draft)
			}
		})
	}
}
