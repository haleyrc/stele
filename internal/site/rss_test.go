package site_test

import (
	"bytes"
	"encoding/xml"
	"strings"
	"testing"

	"github.com/haleyrc/assert"
	"github.com/haleyrc/stele/internal/site"
	"github.com/haleyrc/stele/internal/testutil"
)

func TestRSS_Render(t *testing.T) {
	s := testutil.TestSite()
	rss := site.NewRSSFeed(s)

	var buff bytes.Buffer
	err := rss.Render(&buff)
	assert.OK(t, err).Fatal()

	output := buff.String()

	// Verify XML declaration
	if !strings.HasPrefix(output, `<?xml version="1.0" encoding="UTF-8" ?>`) {
		t.Error("expected RSS output to start with XML declaration")
	}

	// Verify RSS contains expected content using string matching
	// (XML parsing with namespaces is complex, so we use simple checks)
	if !strings.Contains(output, "<title>Alice Codes</title>") {
		t.Error("expected RSS to contain site title")
	}
	if !strings.Contains(output, "<link>https://alice.dev</link>") {
		t.Error("expected RSS to contain site link")
	}
	if !strings.Contains(output, "<description>Alice&#39;s development blog covering Go, web APIs, and software engineering</description>") {
		t.Error("expected RSS to contain site description")
	}
	if !strings.Contains(output, "<language>en</language>") {
		t.Error("expected RSS to contain language")
	}
	if !strings.Contains(output, "<category>technology</category>") {
		t.Error("expected RSS to contain technology category")
	}
	if !strings.Contains(output, "<category>programming</category>") {
		t.Error("expected RSS to contain programming category")
	}
	if !strings.Contains(output, "<category>web-development</category>") {
		t.Error("expected RSS to contain web-development category")
	}

	// Verify copyright format (should be "Copyright YYYY Alice Smith")
	if !strings.Contains(output, "Copyright") || !strings.Contains(output, "Alice Smith") {
		t.Error("expected RSS to contain copyright with author name")
	}

	// Verify lastBuildDate is present (don't check exact value since it's dynamic)
	if !strings.Contains(output, "<lastBuildDate>") {
		t.Error("expected RSS to contain lastBuildDate")
	}

	// Verify items
	if !strings.Contains(output, "<title>Getting Started with Go</title>") {
		t.Error("expected RSS to contain first post title")
	}
	if !strings.Contains(output, "<link>https://alice.dev/posts/getting-started-with-go</link>") {
		t.Error("expected RSS to contain first post link")
	}
	if !strings.Contains(output, "<guid>https://alice.dev/posts/getting-started-with-go</guid>") {
		t.Error("expected RSS to contain first post guid")
	}
	if !strings.Contains(output, "<description>A comprehensive beginner&#39;s guide to the Go programming language</description>") {
		t.Error("expected RSS to contain first post description")
	}
	if !strings.Contains(output, "<pubDate>") {
		t.Error("expected RSS to contain pubDate for posts")
	}
}

func TestSite_RenderRSSFeed(t *testing.T) {
	s := testutil.TestSite()

	var buff bytes.Buffer
	feed := s.RSSFeed()
	err := feed.Render(&buff)
	assert.OK(t, err).Fatal()

	output := buff.String()

	// Verify XML is valid and contains expected content
	var parsed struct {
		XMLName xml.Name `xml:"rss"`
		Channel struct {
			Title string `xml:"title"`
			Items []struct {
				Title string `xml:"title"`
			} `xml:"item"`
		} `xml:"channel"`
	}

	err = xml.Unmarshal([]byte(output), &parsed)
	assert.OK(t, err).Fatal()

	assert.Equal(t, "title", "Alice Codes", parsed.Channel.Title)
	assert.Equal(t, "item count", 4, len(parsed.Channel.Items))
}
