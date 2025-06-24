package template

import (
	"context"
	"fmt"
	"io"

	"github.com/haleyrc/stele/internal/site"
	"github.com/haleyrc/stele/internal/template/pages"
)

// TemplateRenderer implements the site.Renderer interface using basic HTML.
type TemplateRenderer struct{}

// NewTemplateRenderer creates a new template renderer.
func NewTemplateRenderer() *TemplateRenderer {
	return &TemplateRenderer{}
}

// RenderIndex renders the index page using templ components.
func (r *TemplateRenderer) RenderIndex(ctx context.Context, w io.Writer, site *site.Site) error {
	recentPosts := site.Posts.Recent(11)
	latestPost, remainingPosts := recentPosts.Head()
	return Layout("Home", site, pages.Index(latestPost, remainingPosts)).Render(ctx, w)
}

// RenderPost renders a single post page using templ components.
func (r *TemplateRenderer) RenderPost(ctx context.Context, w io.Writer, site *site.Site, post *site.Post) error {
	seriesInfo := site.Series.GetSeriesInfo(post.Slug)
	return Layout(post.Frontmatter.Title, site, pages.Post(post, seriesInfo)).Render(ctx, w)
}

// RenderTagIndex renders the tags index page listing all tags with post counts.
func (r *TemplateRenderer) RenderTagIndex(ctx context.Context, w io.Writer, site *site.Site) error {
	postsByTag := site.Posts.IndexByTag()
	return Layout("Tags", site, pages.PostIndex(postsByTag, "/tags/")).Render(ctx, w)
}

// RenderTagPage renders a page showing all posts with a specific tag.
func (r *TemplateRenderer) RenderTagPage(ctx context.Context, w io.Writer, site *site.Site, tag string, posts site.Posts) error {
	heading := fmt.Sprintf("Posts tagged %q", tag)
	return Layout(heading, site, pages.PostList(heading, posts)).Render(ctx, w)
}

// RenderArchiveIndex renders the archive index page listing all years with post
// counts.
func (r *TemplateRenderer) RenderArchiveIndex(ctx context.Context, w io.Writer, site *site.Site) error {
	postsByYear := site.Posts.IndexByYear()
	return Layout("Archive", site, pages.PostIndex(postsByYear, "/archive/")).Render(ctx, w)
}

// RenderArchivePage renders a page showing all posts from a specific year.
func (r *TemplateRenderer) RenderArchivePage(ctx context.Context, w io.Writer, site *site.Site, year string, posts site.Posts) error {
	heading := fmt.Sprintf("Posts from %s", year)
	return Layout(heading, site, pages.PostList(heading, posts)).Render(ctx, w)
}

// RenderAbout renders the about page using templ components.
func (r *TemplateRenderer) RenderAbout(ctx context.Context, w io.Writer, s *site.Site, about *site.About) error {
	return Layout("About", s, pages.About(s, *about)).Render(ctx, w)
}

// RenderNotesIndex renders the notes index page.
func (r *TemplateRenderer) RenderNotesIndex(ctx context.Context, w io.Writer, s *site.Site) error {
	return Layout("Notes", s, pages.NotesIndex(s)).Render(ctx, w)
}

// RenderNote renders a single note page using templ components.
func (r *TemplateRenderer) RenderNote(ctx context.Context, w io.Writer, s *site.Site, note *site.Note) error {
	return Layout(note.Frontmatter.Title, s, pages.Note(*note)).Render(ctx, w)
}

// RenderNoteTagIndex renders the note tags index page listing all tags with note counts.
func (r *TemplateRenderer) RenderNoteTagIndex(ctx context.Context, w io.Writer, s *site.Site) error {
	notesByTag := s.Notes.IndexByTag()
	return Layout("Note Tags", s, pages.NoteTagIndex(notesByTag)).Render(ctx, w)
}

// RenderNoteTagPage renders a page showing all notes with a specific tag.
func (r *TemplateRenderer) RenderNoteTagPage(ctx context.Context, w io.Writer, s *site.Site, tag string, notes site.Notes) error {
	heading := fmt.Sprintf("Notes tagged %q", tag)
	return Layout(heading, s, pages.NoteTagPage(heading, notes)).Render(ctx, w)
}

// Render404 renders a 404 error page.
func (r *TemplateRenderer) Render404(ctx context.Context, w io.Writer, site *site.Site) error {
	return Layout("404", site, pages.NotFound()).Render(ctx, w)
}

// RenderManifest renders the web manifest as JSON.
func (r *TemplateRenderer) RenderManifest(ctx context.Context, w io.Writer, site *site.Site, manifest *site.Manifest) error {
	return manifest.Render(w)
}

// RenderRSSFeed renders the RSS feed as XML.
func (r *TemplateRenderer) RenderRSSFeed(ctx context.Context, w io.Writer, site *site.Site, feed *site.RSSFeed) error {
	return feed.Render(w)
}

// RenderSeriesIndex renders the index page for a series.
func (r *TemplateRenderer) RenderSeriesIndex(ctx context.Context, w io.Writer, site *site.Site, series *site.Series) error {
	return Layout(series.Metadata.Name, site, pages.SeriesIndex(*series)).Render(ctx, w)
}
