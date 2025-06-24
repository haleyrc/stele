// Package site provides the primary domain types and API for working with
// static blog sites.
package site

import (
	"context"
	"io"
)

// Renderer handles rendering of site content to writers.
type Renderer interface {
	Render404(ctx context.Context, w io.Writer, site *Site) error
	RenderAbout(ctx context.Context, w io.Writer, site *Site, about *About) error
	RenderArchiveIndex(ctx context.Context, w io.Writer, site *Site) error
	RenderArchivePage(ctx context.Context, w io.Writer, site *Site, year string, posts Posts) error
	RenderIndex(ctx context.Context, w io.Writer, site *Site) error
	RenderManifest(ctx context.Context, w io.Writer, site *Site, manifest *Manifest) error
	RenderNote(ctx context.Context, w io.Writer, site *Site, note *Note) error
	RenderNotesIndex(ctx context.Context, w io.Writer, site *Site) error
	RenderNoteTagIndex(ctx context.Context, w io.Writer, site *Site) error
	RenderNoteTagPage(ctx context.Context, w io.Writer, site *Site, tag string, notes Notes) error
	RenderPost(ctx context.Context, w io.Writer, site *Site, post *Post) error
	RenderRSSFeed(ctx context.Context, w io.Writer, site *Site, feed *RSSFeed) error
	RenderSeriesIndex(ctx context.Context, w io.Writer, site *Site, series *Series) error
	RenderTagIndex(ctx context.Context, w io.Writer, site *Site) error
	RenderTagPage(ctx context.Context, w io.Writer, site *Site, tag string, posts Posts) error
}
