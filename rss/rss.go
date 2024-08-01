package rss

import (
	"fmt"
	"time"

	"github.com/haleyrc/stele/site"
)

type AtomLink struct {
	Href string `xml:"href,attr"`
	Rel  string `xml:"rel,attr"`
	Type string `xml:"type,attr"`
}

type builder struct {
	buildTime time.Time
}

type BuildOption func(*builder)

func WithBuildTime(t time.Time) BuildOption {
	return func(b *builder) {
		b.buildTime = t
	}
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
	Version string  `xml:"version,attr"`
	NSAtom  string  `xml:"xmlns:atom,attr"`
	Channel Channel `xml:"channel"`
}

func Build(s *site.Site, opts ...BuildOption) (*Feed, error) {
	b := builder{
		buildTime: time.Now(),
	}

	for _, opt := range opts {
		opt(&b)
	}

	feed := &Feed{
		Version: "2.0",
		NSAtom:  "http://www.w3.org/2005/Atom",
		Channel: Channel{
			AtomLink: AtomLink{
				Href: s.Config.BaseURL + "/rss.xml",
				Rel:  "self",
				Type: "application/rss+xml",
			},
			Category: s.Config.Categories,
			Copyright: fmt.Sprintf(
				"Copyright %d %s",
				s.Index.Posts.First().Timestamp.Year(),
				s.Config.Author,
			),
			Description:   s.Config.Description,
			Language:      "en",
			LastBuildDate: b.buildTime.Format(time.RFC1123Z),
			Link:          s.Config.BaseURL,
			Title:         s.Config.Name,
		},
	}

	for _, p := range s.Index.Posts {
		feed.Channel.Items = append(feed.Channel.Items, Item{
			Title:       p.Title,
			Link:        fmt.Sprintf("%s/posts/%s", s.Config.BaseURL, p.Slug),
			GUID:        fmt.Sprintf("%s/posts/%s", s.Config.BaseURL, p.Slug),
			Description: p.Description,
			Category:    p.Tags,
			PubDate:     p.Timestamp.Format(time.RFC1123Z),
		})
	}

	return feed, nil
}

type Item struct {
	Title       string   `xml:"title"`
	Link        string   `xml:"link"`
	GUID        string   `xml:"guid"`
	Description string   `xml:"description,omitempty"`
	Category    []string `xml:"category"`
	PubDate     string   `xml:"pubDate"`
}
