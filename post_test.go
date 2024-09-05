package stele_test

import (
	"path/filepath"
	"slices"
	"testing"

	"github.com/haleyrc/stele"
)

func TestNewPost(t *testing.T) {
	path := filepath.Join("testdata", "posts", "20220103-first-post.md")
	post, err := stele.NewPost(path)
	if err != nil {
		t.Fatal(err)
	}

	if want := "First Post"; post.Title != want {
		t.Errorf("expected post.Title = %q, but it was %q", want, post.Title)
	}
	if want := "The first post"; post.Description != want {
		t.Errorf("expected post.Description = %q, but it was %q", want, post.Description)
	}
	if want := []string{"go", "react"}; !slices.Equal(want, post.Tags) {
		t.Errorf("expected post.Tags = %v, but it was %v", want, post.Tags)
	}
}
