package rss_test

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"

	"github.com/haleyrc/stele/rss"
	"github.com/haleyrc/stele/site"
	"github.com/haleyrc/stele/site/config"
	"github.com/haleyrc/stele/site/index"
)

func TestBuild(t *testing.T) {
	site := &site.Site{
		Config: &config.Config{
			Author:      "John Smith",
			BaseURL:     "https://example.com",
			Categories:  []string{"programming", "blog"},
			Description: "This is a test blog",
			Name:        "Test",
		},
		Index: &index.Index{
			Posts: index.Posts{
				{
					Description: "This post is too new",
					Slug:        "newest-post",
					Tags:        []string{"go", "react"},
					Timestamp:   time.Unix(1722523275, 0),
					Title:       "Newest Post",
				},
				{
					Description: "This post is too old",
					Slug:        "oldest-post",
					Tags:        []string{"go", "react"},
					Timestamp:   time.Unix(1722523275, 0).AddDate(-3, 0, 0),
					Title:       "Oldest Post",
				},
				{
					Description: "This post is just right",
					Slug:        "middle-post",
					Tags:        []string{"go", "react"},
					Timestamp:   time.Unix(1722523275, 0).AddDate(-1, 0, 0),
					Title:       "Middle Post",
				},
			},
		},
	}
	want := &rss.Feed{
		Version: "2.0",
		NSAtom:  "http://www.w3.org/2005/Atom",
		Channel: rss.Channel{
			AtomLink: rss.AtomLink{
				Href: "https://example.com/rss.xml",
				Rel:  "self",
				Type: "application/rss+xml",
			},
			Category:    []string{"programming", "blog"},
			Copyright:   "Copyright 2021 John Smith",
			Description: "This is a test blog",
			Items: []rss.Item{
				{
					Title:       "Newest Post",
					Link:        "https://example.com/posts/newest-post",
					GUID:        "https://example.com/posts/newest-post",
					Description: "This post is too new",
					Category:    []string{"go", "react"},
					PubDate:     "Thu, 01 Aug 2024 10:41:15 -0400",
				},
				{
					Title:       "Oldest Post",
					Link:        "https://example.com/posts/oldest-post",
					GUID:        "https://example.com/posts/oldest-post",
					Description: "This post is too old",
					Category:    []string{"go", "react"},
					PubDate:     "Sun, 01 Aug 2021 10:41:15 -0400",
				},
				{
					Title:       "Middle Post",
					Link:        "https://example.com/posts/middle-post",
					GUID:        "https://example.com/posts/middle-post",
					Description: "This post is just right",
					Category:    []string{"go", "react"},
					PubDate:     "Tue, 01 Aug 2023 10:41:15 -0400",
				},
			},
			Language:      "en",
			LastBuildDate: "Thu, 01 Aug 2024 10:41:15 -0400",
			Link:          "https://example.com",
			Title:         "Test",
		},
	}

	got, err := rss.Build(site,
		rss.WithBuildTime(time.Unix(1722523275, 0)),
	)
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("incorrect feed (-want, +got):\n%s", diff)
	}
}
