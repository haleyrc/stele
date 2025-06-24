package site

import (
	"encoding/xml"
	"fmt"
	"io"
	"time"
)

// RSSFeed represents an RSS feed for a site.
type RSSFeed struct {
	XMLName xml.Name `xml:"rss"`

	// The RSS version number (e.g., "2.0").
	Version string `xml:"version,attr"`

	// The Atom namespace for atom:link elements.
	NSAtom string `xml:"xmlns:atom,attr"`

	// The feed metadata and items.
	Channel RSSFeedChannel `xml:"channel"`
}

// RSSFeedChannel represents the RSS channel element.
type RSSFeedChannel struct {
	// The name of the channel (feed title).
	Title string `xml:"title"`

	// The URL to the HTML website corresponding to the channel.
	Link string `xml:"link"`

	// Autodiscovery information for feed readers.
	RSSFeedChannelAtomLink RSSFeedChannelAtomLink `xml:"atom:link"`

	// A brief description of the channel.
	Description string `xml:"description"`

	// One or more categories that the channel belongs to.
	Category []string `xml:"category"`

	// The copyright notice for content in the channel.
	Copyright string `xml:"copyright"`

	// The language the channel is written in.
	Language string `xml:"language"`

	// When the content of the channel last changed.
	LastBuildDate string `xml:"lastBuildDate"`

	// The individual posts/articles in the feed.
	Items []RSSFeedChannelItem `xml:"item"`
}

// RSSFeedChannelAtomLink represents the atom:link element for RSS.
type RSSFeedChannelAtomLink struct {
	// The URL of the resource.
	Href string `xml:"href,attr"`

	// The relationship between the current document and the linked resource.
	Rel string `xml:"rel,attr"`

	// The media type of the linked resource (e.g., "application/rss+xml",
	// "text/html", "application/atom+xml").
	Type string `xml:"type,attr"`
}

// RSSFeedChannelItem represents an individual RSS item.
type RSSFeedChannelItem struct {
	// The title of the item.
	Title string `xml:"title"`

	// The URL of the item.
	Link string `xml:"link"`

	// A globally unique identifier for the item.
	GUID string `xml:"guid"`

	// A synopsis of the item content.
	Description string `xml:"description,omitempty"`

	// One or more categories that the item belongs to.
	Category []string `xml:"category"`

	// When the item was published.
	PubDate string `xml:"pubDate"`
}

// NewRSSFeed creates a new RSS feed for the given site.
func NewRSSFeed(s *Site) *RSSFeed {
	rss := &RSSFeed{
		Version: "2.0",
		NSAtom:  "http://www.w3.org/2005/Atom",
		Channel: RSSFeedChannel{
			Title: s.Config.Title,
			Link:  s.Config.BaseURL,
			RSSFeedChannelAtomLink: RSSFeedChannelAtomLink{
				Href: s.Config.BaseURL + "/rss.xml",
				Rel:  "self",
				Type: "application/rss+xml",
			},
			Description: s.Config.Description,
			Category:    s.Config.Categories,
			Copyright: fmt.Sprintf(
				"Copyright %d %s",
				s.CopyrightYear(),
				s.Config.Author,
			),
			Language:      "en",
			LastBuildDate: time.Now().Format(time.RFC1123),
			Items:         []RSSFeedChannelItem{},
		},
	}

	for _, post := range s.Posts {
		item := RSSFeedChannelItem{
			Title:       post.Frontmatter.Title,
			Link:        fmt.Sprintf("%s/posts/%s", s.Config.BaseURL, post.Slug),
			GUID:        fmt.Sprintf("%s/posts/%s", s.Config.BaseURL, post.Slug),
			Description: post.Frontmatter.Description,
			Category:    post.Frontmatter.Tags,
			PubDate:     post.Frontmatter.Timestamp.Format(time.RFC1123),
		}
		rss.Channel.Items = append(rss.Channel.Items, item)
	}

	return rss
}

// Render writes the RSS feed as XML to the provided writer.
func (r *RSSFeed) Render(w io.Writer) error {
	if _, err := fmt.Fprintln(w, `<?xml version="1.0" encoding="UTF-8" ?>`); err != nil {
		return fmt.Errorf("rss: render: %w", err)
	}

	bytes, err := xml.MarshalIndent(r, "", "  ")
	if err != nil {
		return fmt.Errorf("rss: render: %w", err)
	}

	if _, err := w.Write(bytes); err != nil {
		return fmt.Errorf("rss: render: %w", err)
	}

	return nil
}
