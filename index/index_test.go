package index_test

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"

	"github.com/haleyrc/stele/index"
)

func TestNew(t *testing.T) {
	want := &index.Index{
		Pages: []index.Page{
			{
				Path: "testdata/pages/about.html",
				Slug: "about",
			},
			{
				Path: "testdata/pages/contact-us.html",
				Slug: "contact-us",
			},
		},
		Posts: index.Posts{
			{
				Description: "The first post",
				Path:        "testdata/posts/20220103-first-post.md",
				Slug:        "first-post",
				Tags:        []string{"go", "react"},
				Timestamp:   timestamp("20220103"),
				Title:       "First Post",
			},
			{
				Description: "The second post",
				Path:        "testdata/posts/20240406-second-post.md",
				Slug:        "second-post",
				Tags:        []string{"go", "react"},
				Timestamp:   timestamp("20240406"),
				Title:       "Second Post",
			},
		},
	}

	idx, err := index.New("testdata")
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(want, idx); diff != "" {
		t.Fatalf("incorrect index (-want, +got):\n%s", diff)
	}
}

func timestamp(s string) time.Time {
	ts, _ := time.Parse("20060102", s)
	return ts
}
