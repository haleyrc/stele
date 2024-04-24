package stele

import (
	"bytes"
	"context"
	_ "embed"
	"testing"

	"github.com/google/go-cmp/cmp"
)

//go:embed testdata/rss.xml
var want string

func TestRSS(t *testing.T) {
	c := channel{
		Title:       "Test Channel",
		Link:        "https://example.com",
		Description: "An example RSS feed for testing",
		Category:    []string{"Personal blog"},
		Copyright:   "Copyright 2022 Ryan Haley",
		Image: &image{
			Link:        "https://example.com",
			Title:       "Test Channel",
			URL:         "https://example.com/masthead.png",
			Description: "An example RSS feed for testing",
			Height:      32,
			Width:       96,
		},
		Language:      "en",
		LastBuildDate: "Thur, 18 Apr 2024 20:39:44 EST",
		Items: []item{{
			Title:       "The latest post",
			Link:        "https://example.com/posts/the-latest-post",
			Description: "The latest post on a fake blog",
			Category:    []string{"go", "react"},
			PubDate:     "Fri, 05 Oct 2007 09:00:00 EST",
		}},
	}

	var buff bytes.Buffer
	if err := c.Render(context.Background(), &buff); err != nil {
		t.Fatal("unexpected error:", err)
	}

	got := buff.String()
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Expected values to be equal:\n%s", diff)
	}
}
