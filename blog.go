package stele

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strconv"
	"strings"
	"time"

	katex "github.com/FurqanSoftware/goldmark-katex"
	"github.com/a-h/templ"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"go.abhg.dev/goldmark/frontmatter"
	"golang.org/x/exp/maps"
	"gopkg.in/yaml.v3"
)

func Build(ctx context.Context, b *Blog, t Theme) error {
	if err := b.Load(); err != nil {
		return errf("build", err)
	}
	if err := clean(); err != nil {
		return errf("build", err)
	}
	if err := makeDir(); err != nil {
		return errf("build", err)
	}
	if err := buildArchive(ctx, b, t); err != nil {
		return errf("build", err)
	}
	if err := buildTags(ctx, b, t); err != nil {
		return errf("build", err)
	}
	if err := buildIndex(ctx, b, t); err != nil {
		return errf("build", err)
	}
	if err := buildPages(ctx, b, t); err != nil {
		return errf("build", err)
	}
	if err := buildPosts(ctx, b, t); err != nil {
		return errf("build", err)
	}
	if err := buildRSS(ctx, b); err != nil {
		return errf("build", err)
	}
	if err := buildManifest(ctx, b); err != nil {
		return errf("build", err)
	}
	return nil
}

type Blog struct {
	Config Config
	Pages  Pages
	Posts  Posts
}

func New() (*Blog, error) {
	cfg, err := LoadConfig(filepath.Join(".", "config.yml"))
	if err != nil {
		return nil, errf("new", err)
	}
	return &Blog{Config: *cfg}, nil
}

func (b *Blog) Load() error {
	var err error

	b.Pages, err = LoadPages()
	if err != nil {
		return errf("blog: load", err)
	}

	b.Posts, err = LoadPosts()
	if err != nil {
		return errf("blog: load", err)
	}

	return nil
}

type atomLink struct {
	Href string `xml:"href,attr"`
	Rel  string `xml:"rel,attr"`
	Type string `xml:"type,attr"`
}

type channel struct {
	Title         string   `xml:"title"`
	Link          string   `xml:"link"`
	AtomLink      atomLink `xml:"atom:link"`
	Description   string   `xml:"description"`
	Category      []string `xml:"category"`
	Copyright     string   `xml:"copyright"`
	Image         *image   `xml:"image,omitempty"`
	Language      string   `xml:"language"`
	LastBuildDate string   `xml:"lastBuildDate"`
	Items         []item   `xml:"item"`
}

func blogToChannel(b *Blog) channel {
	c := channel{
		Title: b.Config.Name,
		Link:  b.Config.BaseURL,
		AtomLink: atomLink{
			Href: b.Config.BaseURL + "/rss.xml",
			Rel:  "self",
			Type: "application/rss+xml",
		},
		Description: b.Config.Description,
		Category:    b.Config.Categories,
		Copyright: fmt.Sprintf(
			"Copyright %d %s",
			b.Posts.Latest().Timestamp.Year(),
			b.Config.Author,
		),
		Language:      "en", // TODO
		LastBuildDate: time.Now().Format(time.RFC822),
	}

	for _, p := range b.Posts {
		c.Items = append(c.Items, item{
			Title:       p.Title,
			Link:        fmt.Sprintf("%s/posts/%s", b.Config.BaseURL, p.Slug),
			Description: p.Description,
			Category:    p.Tags,
			PubDate:     p.Timestamp.Format(time.RFC822),
		})
	}

	return c
}

func (c channel) Render(ctx context.Context, w io.Writer) error {
	fmt.Fprintln(w, "<?xml version=\"1.0\" encoding=\"UTF-8\" ?>")
	fmt.Fprintln(w, "<rss version=\"2.0\" xmlns:atom=\"http://www.w3.org/2005/Atom\">")
	fmt.Fprintln(w)

	enc := xml.NewEncoder(w)
	enc.Indent("", "  ")
	if err := enc.Encode(c); err != nil {
		return errf("channel: render", err)
	}

	fmt.Fprintln(w)

	return nil
}

type Config struct {
	Author      string   `yaml:"author"`
	BaseURL     string   `yaml:"baseURL"`
	Categories  []string `yaml:"categories"`
	Description string   `yaml:"description"`
	Menu        []Link   `yaml:"menu"`
	Name        string   `yaml:"name"`
}

func LoadConfig(path string) (*Config, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, errf("load config", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(bytes, &cfg); err != nil {
		return nil, errf("load config", err)
	}

	return &cfg, nil
}

type icon struct {
	Purpose string `json:"purpose"`
	Sizes   string `json:"sizes"`
	Source  string `json:"src"`
	Type    string `json:"type"`
}

type image struct {
	Link        string `xml:"link"`
	Title       string `xml:"title"`
	URL         string `xml:"url"`
	Description string `xml:"description"`
	Height      int    `xml:"height"`
	Width       int    `xml:"width"`
}

type item struct {
	Title       string   `xml:"title"`
	Link        string   `xml:"link"`
	Description string   `xml:"description,omitempty"`
	Category    []string `xml:"category"`
	PubDate     string   `xml:"pubDate"`
}

type Link struct {
	Label string `yaml:"label"`
	Path  string `yaml:"path"`
}

type manifest struct {
	BackgroundColor string   `json:"background_color"`
	Categories      []string `json:"categories"`
	Description     string   `json:"description"`
	Display         string   `json:"display"`
	Icons           []icon   `json:"icons"`
	Name            string   `json:"name"`
	StartURL        string   `json:"start_url"`
}

func blogToManifest(b *Blog) manifest {
	m := manifest{
		BackgroundColor: "white",
		Categories:      b.Config.Categories,
		Description:     b.Config.Description,
		Display:         "fullscreen",
		Icons:           []icon{},
		Name:            b.Config.Name,
		StartURL:        b.Config.BaseURL,
	}
	return m
}

func (m manifest) Render(ctx context.Context, w io.Writer) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	if err := enc.Encode(m); err != nil {
		return errf("manifest: render", err)
	}
	return nil
}

type Page struct {
	Content string
	Slug    string
}

type Pages []Page

func LoadPages() (Pages, error) {
	files, err := filepath.Glob(filepath.Join(".", "pages", "*.html"))
	if err != nil {
		return nil, errf("load pages", err)
	}

	pages := make([]Page, 0, len(files))
	for _, file := range files {
		contents, err := os.ReadFile(file)
		if err != nil {
			return nil, errf("load pages", err)
		}

		pages = append(pages, Page{
			Content: string(contents),
			Slug:    strings.TrimSuffix(filepath.Base(file), ".html"),
		})
	}

	return pages, nil
}

type Post struct {
	Content     string
	Description string
	Slug        string
	Tags        []string
	Timestamp   time.Time
	Title       string
}

type PostBucket struct {
	Key   string
	Posts []Post
}

type Posts []Post

func LoadPosts() (Posts, error) {
	md := goldmark.New(goldmark.WithExtensions(
		extension.GFM,
		&frontmatter.Extender{},
		&katex.Extender{},
	))

	files, err := filepath.Glob(filepath.Join(".", "posts", "*.md"))
	if err != nil {
		return nil, errf("load posts", err)
	}

	posts := make([]Post, 0, len(files))
	for _, file := range files {
		contents, err := os.ReadFile(file)
		if err != nil {
			return nil, errf("load posts", err)
		}

		var buff bytes.Buffer
		ctx := parser.NewContext()
		if err := md.Convert(contents, &buff, parser.WithContext(ctx)); err != nil {
			return nil, errf("load posts", err)
		}

		var meta struct {
			Description string   `yaml:"description"`
			Tags        []string `yaml:"tags"`
			Title       string   `yaml:"title"`
		}
		if err := frontmatter.Get(ctx).Decode(&meta); err != nil {
			return nil, errf("load posts", err)
		}

		slug, timestamp, err := parsePostName(file)
		if err != nil {
			return nil, errf("load posts", err)
		}

		posts = append(posts, Post{
			Content:     buff.String(),
			Description: meta.Description,
			Slug:        slug,
			Tags:        meta.Tags,
			Timestamp:   timestamp,
			Title:       meta.Title,
		})
	}

	sort.Slice(posts, func(i, j int) bool {
		return posts[j].Timestamp.Before(posts[i].Timestamp)
	})

	return posts, nil
}

func (ps Posts) ByTag() []PostBucket {
	pm := map[string][]Post{}
	for _, p := range ps {
		for _, t := range p.Tags {
			if _, found := pm[t]; !found {
				pm[t] = []Post{}
			}
			pm[t] = append(pm[t], p)
		}
	}

	keys := maps.Keys(pm)
	sort.Strings(keys)
	buckets := make([]PostBucket, 0, len(keys))
	for _, key := range keys {
		buckets = append(buckets, PostBucket{
			Key:   key,
			Posts: pm[key],
		})
	}

	return buckets
}

func (ps Posts) ByYear() []PostBucket {
	pm := map[string][]Post{}
	for _, p := range ps {
		year := strconv.Itoa(p.Timestamp.Year())
		if _, found := pm[year]; !found {
			pm[year] = []Post{}
		}
		pm[year] = append(pm[year], p)
	}

	keys := maps.Keys(pm)
	sort.Strings(keys)
	slices.Reverse(keys)
	buckets := make([]PostBucket, 0, len(keys))
	for _, key := range keys {
		buckets = append(buckets, PostBucket{
			Key:   key,
			Posts: pm[key],
		})
	}

	return buckets
}

func (ps Posts) First() Post {
	var first *Post
	for _, p := range ps {
		if first == nil || p.Timestamp.Before(first.Timestamp) {
			first = &p
		}
	}
	return *first
}

func (ps Posts) Latest() Post {
	var latest *Post
	for _, p := range ps {
		if latest == nil || p.Timestamp.After(latest.Timestamp) {
			latest = &p
		}
	}
	return *latest
}

func (ps Posts) MostRecent(max int) []Post {
	if len(ps) < max {
		max = len(ps)
	}
	return ps[:max]
}

type Renderable interface {
	Render(ctx context.Context, w io.Writer) error
}

type Theme interface {
	Archive(*Blog) templ.Component
	Index(*Blog) templ.Component
	Page(*Blog, Page) templ.Component
	Post(*Blog, Post) templ.Component
	PostsForTag(*Blog, PostBucket) templ.Component
	PostsForYear(*Blog, PostBucket) templ.Component
	Tags(*Blog) templ.Component
}

func buildArchive(ctx context.Context, b *Blog, t Theme) error {
	path := filepath.Join(".", "public", "archive.html")
	if err := render(ctx, path, t.Archive(b)); err != nil {
		return errf("build archive", err)
	}
	for _, bucket := range b.Posts.ByYear() {
		path := filepath.Join(".", "public", "archive", bucket.Key+".html")
		if err := render(ctx, path, t.PostsForYear(b, bucket)); err != nil {
			return errf("build archive", err)
		}
	}
	return nil
}

func buildIndex(ctx context.Context, b *Blog, t Theme) error {
	path := filepath.Join(".", "public", "index.html")
	if err := render(ctx, path, t.Index(b)); err != nil {
		return errf("build index", err)
	}
	return nil
}

func buildManifest(ctx context.Context, b *Blog) error {
	m := blogToManifest(b)
	path := filepath.Join(".", "public", "manifest.webmanifest")
	if err := render(ctx, path, m); err != nil {
		return errf("build manifest", err)
	}
	return nil
}

func buildPages(ctx context.Context, b *Blog, t Theme) error {
	for _, p := range b.Pages {
		path := filepath.Join(".", "public", p.Slug+".html")
		if err := render(ctx, path, t.Page(b, p)); err != nil {
			return errf("build pages", err)
		}
	}
	return nil
}

func buildPosts(ctx context.Context, b *Blog, t Theme) error {
	for _, p := range b.Posts {
		path := filepath.Join(".", "public", "posts", p.Slug+".html")
		if err := render(ctx, path, t.Post(b, p)); err != nil {
			return errf("build posts", err)
		}
	}
	return nil
}

func buildRSS(ctx context.Context, b *Blog) error {
	c := blogToChannel(b)
	path := filepath.Join(".", "public", "rss.xml")
	if err := render(ctx, path, c); err != nil {
		return errf("build rss", err)
	}
	return nil
}

func buildTags(ctx context.Context, b *Blog, t Theme) error {
	path := filepath.Join(".", "public", "tags.html")
	if err := render(ctx, path, t.Tags(b)); err != nil {
		return errf("build tags", err)
	}
	for _, bucket := range b.Posts.ByTag() {
		path := filepath.Join(".", "public", "tags", bucket.Key+".html")
		if err := render(ctx, path, t.PostsForTag(b, bucket)); err != nil {
			return errf("build tags", err)
		}
	}
	return nil
}

func clean() error {
	path := filepath.Join(".", "public")
	if err := os.RemoveAll(path); err != nil {
		return errf("clean", err)
	}
	return nil
}

func errf(op string, err error) error {
	return fmt.Errorf("%s: %w", op, err)
}

func makeDir() error {
	dirs := [][]string{
		{".", "public"},
		{".", "public", "posts"},
		{".", "public", "archive"},
		{".", "public", "tags"},
	}
	for _, parts := range dirs {
		path := filepath.Join(parts...)
		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			return errf("make dir", err)
		}
	}
	return nil
}

func parsePostName(filename string) (string, time.Time, error) {
	name := strings.TrimSuffix(filepath.Base(filename), ".md")
	nameParts := strings.SplitN(name, "-", 2)

	ts, err := time.Parse("20060102", nameParts[0])
	if err != nil {
		return "", time.Time{}, errf("parse post name", err)
	}

	return nameParts[1], ts, nil
}

func render(ctx context.Context, path string, r Renderable) error {
	var buff bytes.Buffer
	if err := r.Render(ctx, &buff); err != nil {
		return errf("render", err)
	}
	if err := os.WriteFile(path, buff.Bytes(), os.ModePerm); err != nil {
		return errf("render", err)
	}
	return nil
}
