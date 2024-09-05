package stele

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"time"
)

type atomLink struct {
	Href string `xml:"href,attr"`
	Rel  string `xml:"rel,attr"`
	Type string `xml:"type,attr"`
}

type channel struct {
	Title       string   `xml:"title"`
	Link        string   `xml:"link"`
	AtomLink    atomLink `xml:"atom:link"`
	Description string   `xml:"description"`
	Category    []string `xml:"category"`
	Copyright   string   `xml:"copyright"`
	// Image         *image   `xml:"image,omitempty"`
	Language      string `xml:"language"`
	LastBuildDate string `xml:"lastBuildDate"`
	Items         []item `xml:"item"`
}

type feed struct {
	XMLName xml.Name `xml:"rss"`
	Version string   `xml:"version,attr"`
	NSAtom  string   `xml:"xmlns:atom,attr"`
	Channel channel  `xml:"channel"`
}

func newFeed(cfg *Config, posts Posts) (*feed, error) {
	feed := &feed{
		Version: "2.0",
		NSAtom:  "http://www.w3.org/2005/Atom",
		Channel: channel{
			AtomLink: atomLink{
				Href: cfg.BaseURL + "/rss.xml",
				Rel:  "self",
				Type: "application/rss+xml",
			},
			Category:      cfg.Categories,
			Description:   cfg.Description,
			Language:      "en",
			LastBuildDate: time.Now().Format(time.RFC1123Z),
			Link:          cfg.BaseURL,
			Title:         cfg.Title,
		},
	}

	if count := len(posts); count > 0 {
		feed.Channel.Copyright = fmt.Sprintf(
			"Copyright %d %s",
			posts[count-1].Timestamp.Year(),
			cfg.Author,
		)

		postsPrefix := fmt.Sprintf("%s/posts", cfg.BaseURL)
		for _, p := range posts {
			item, err := newItem(postsPrefix, p)
			if err != nil {
				return nil, fmt.Errorf("blog: new feed: %w", err)
			}
			feed.Channel.Items = append(feed.Channel.Items, *item)
		}
	}

	return feed, nil
}

type item struct {
	Title       string   `xml:"title"`
	Link        string   `xml:"link"`
	GUID        string   `xml:"guid"`
	Description string   `xml:"description,omitempty"`
	Category    []string `xml:"category"`
	PubDate     string   `xml:"pubDate"`
}

func newItem(baseURL string, p Post) (*item, error) {
	item := &item{
		Title:       p.Title,
		Link:        fmt.Sprintf("%s/%s", baseURL, p.Slug),
		GUID:        fmt.Sprintf("%s/%s", baseURL, p.Slug),
		Description: p.Description,
		Category:    p.Tags,
		PubDate:     p.Timestamp.Format(time.RFC1123Z),
	}
	return item, nil
}

// RenderRSSFeed renders an RSS feed to w based on the configuration and list of
// posts provided.
func RenderRSSFeed(ctx context.Context, w io.Writer, cfg *Config, posts Posts) error {
	feed, err := newFeed(cfg, posts)
	if err != nil {
		return fmt.Errorf("render rss feed: %w", err)
	}

	if _, err := w.Write([]byte("<?xml version=\"1.0\" encoding=\"UTF-8\" ?>\n\n")); err != nil {
		return fmt.Errorf("render rss feed: %w", err)
	}

	bytes, err := xml.MarshalIndent(feed, "", "  ")
	if err != nil {
		return fmt.Errorf("render rss feed: %w", err)
	}

	if _, err := w.Write(bytes); err != nil {
		return fmt.Errorf("render rss feed: %w", err)
	}

	return nil
}
