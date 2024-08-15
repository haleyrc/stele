package blog

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"time"
)

type AtomLink struct {
	Href string `xml:"href,attr"`
	Rel  string `xml:"rel,attr"`
	Type string `xml:"type,attr"`
}

type Channel struct {
	Title       string   `xml:"title"`
	Link        string   `xml:"link"`
	AtomLink    AtomLink `xml:"atom:link"`
	Description string   `xml:"description"`
	Category    []string `xml:"category"`
	Copyright   string   `xml:"copyright"`
	// Image         *image   `xml:"image,omitempty"`
	Language      string `xml:"language"`
	LastBuildDate string `xml:"lastBuildDate"`
	Items         []Item `xml:"item"`
}

type Feed struct {
	XMLName xml.Name `xml:"rss"`
	Version string   `xml:"version,attr"`
	NSAtom  string   `xml:"xmlns:atom,attr"`
	Channel Channel  `xml:"channel"`
}

func NewFeed(cfg *Config, posts Posts) (*Feed, error) {
	feed := &Feed{
		Version: "2.0",
		NSAtom:  "http://www.w3.org/2005/Atom",
		Channel: Channel{
			AtomLink: AtomLink{
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
			item, err := NewItem(postsPrefix, p)
			if err != nil {
				return nil, fmt.Errorf("blog: new feed: %w", err)
			}
			feed.Channel.Items = append(feed.Channel.Items, *item)
		}
	}

	return feed, nil
}

func (f *Feed) Render(ctx context.Context, w io.Writer) error {
	if _, err := w.Write([]byte("<?xml version=\"1.0\" encoding=\"UTF-8\" ?>\n\n")); err != nil {
		return fmt.Errorf("feed: render: %w", err)
	}

	bytes, err := xml.MarshalIndent(f, "", "  ")
	if err != nil {
		return fmt.Errorf("feed: render: %w", err)
	}

	if _, err := w.Write(bytes); err != nil {
		return fmt.Errorf("feed: render: %w", err)
	}

	return nil
}

type Item struct {
	Title       string   `xml:"title"`
	Link        string   `xml:"link"`
	GUID        string   `xml:"guid"`
	Description string   `xml:"description,omitempty"`
	Category    []string `xml:"category"`
	PubDate     string   `xml:"pubDate"`
}

func NewItem(baseURL string, p Post) (*Item, error) {
	item := &Item{
		Title:       p.Title,
		Link:        fmt.Sprintf("%s/%s", baseURL, p.Slug),
		GUID:        fmt.Sprintf("%s/%s", baseURL, p.Slug),
		Description: p.Description,
		Category:    p.Tags,
		PubDate:     p.Timestamp.Format(time.RFC1123Z),
	}
	return item, nil
}
