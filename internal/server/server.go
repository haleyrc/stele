// Package server provides a development HTTP server for previewing site
// content.
package server

import (
	"bytes"
	"context"
	"io"
	"log"
	"net/http"

	"github.com/haleyrc/stele/internal/site"
)

// Server provides HTTP handlers for serving site content.
// It has no knowledge of reloading, caching, or file watching.
type Server struct {
	*http.ServeMux

	// The renderer used to generate HTML output for site pages.
	Renderer site.Renderer
}

// NewServer creates a new content server with all routes registered.
func NewServer(renderer site.Renderer) *Server {
	s := &Server{
		ServeMux: http.NewServeMux(),
		Renderer: renderer,
	}

	// Content endpoints
	s.HandleFunc("GET /", s.HandleIndex)
	s.HandleFunc("GET /about", s.HandleAbout)
	s.HandleFunc("GET /favicon.ico", s.HandleFavicon)
	s.HandleFunc("GET /manifest.webmanifest", s.HandleManifest)
	s.HandleFunc("GET /rss.xml", s.HandleRSS)
	s.HandleFunc("GET /notes", s.HandleNotesIndex)
	s.HandleFunc("GET /notes/{slug}", s.HandleNote)
	s.HandleFunc("GET /notes/tags", s.HandleNoteTagIndex)
	s.HandleFunc("GET /notes/tags/{tag}", s.HandleNoteTagPage)
	s.HandleFunc("GET /posts/{slug}", s.HandlePost)
	s.HandleFunc("GET /posts/{seriesSlug}/{postSlug}", s.HandleSeriesPost)
	s.HandleFunc("GET /tags", s.HandleTagIndex)
	s.HandleFunc("GET /tags/{tag}", s.HandleTagPage)
	s.HandleFunc("GET /archive", s.HandleArchiveIndex)
	s.HandleFunc("GET /archive/{year}", s.HandleArchivePage)
	s.HandleFunc("GET /{seriesSlug}", s.HandleSeriesIndex)

	return s
}

// renderHTML renders HTML content using a buffered approach to ensure errors
// are caught before sending any response to the client.
func (s *Server) renderHTML(w http.ResponseWriter, r *http.Request, handlerName string, renderFn func(context.Context, io.Writer) error) {
	ctx := r.Context()
	var buf bytes.Buffer
	if err := renderFn(ctx, &buf); err != nil {
		log.Printf("ERR: %s: %s: %v", handlerName, r.URL.Path, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(buf.Bytes()) // #nosec G104 - Write errors cannot be handled after headers sent
}

// HandleIndex serves the index page for the site.
func (s *Server) HandleIndex(w http.ResponseWriter, r *http.Request) {
	site := SiteFromContext(r.Context())
	s.renderHTML(w, r, "HandleIndex", func(ctx context.Context, w io.Writer) error {
		return s.Renderer.RenderIndex(ctx, w, site)
	})
}

// HandleAbout serves the about page for the site.
func (s *Server) HandleAbout(w http.ResponseWriter, r *http.Request) {
	site := SiteFromContext(r.Context())
	if site.About == nil {
		s.Handle404(w, r)
		return
	}

	s.renderHTML(w, r, "HandleAbout", func(ctx context.Context, w io.Writer) error {
		return s.Renderer.RenderAbout(ctx, w, site, site.About)
	})
}

// HandleFavicon handles requests for /favicon.ico by returning 204 No Content.
func (s *Server) HandleFavicon(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

// HandleNotesIndex serves the notes index page.
func (s *Server) HandleNotesIndex(w http.ResponseWriter, r *http.Request) {
	site := SiteFromContext(r.Context())
	s.renderHTML(w, r, "HandleNotesIndex", func(ctx context.Context, w io.Writer) error {
		return s.Renderer.RenderNotesIndex(ctx, w, site)
	})
}

// HandleNote serves a single note page.
func (s *Server) HandleNote(w http.ResponseWriter, r *http.Request) {
	site := SiteFromContext(r.Context())
	note := site.Notes.GetBySlug(r.PathValue("slug"))
	if note == nil {
		s.Handle404(w, r)
		return
	}

	s.renderHTML(w, r, "HandleNote", func(ctx context.Context, w io.Writer) error {
		return s.Renderer.RenderNote(ctx, w, site, note)
	})
}

// HandleNoteTagIndex serves the note tags index page.
func (s *Server) HandleNoteTagIndex(w http.ResponseWriter, r *http.Request) {
	site := SiteFromContext(r.Context())
	s.renderHTML(w, r, "HandleNoteTagIndex", func(ctx context.Context, w io.Writer) error {
		return s.Renderer.RenderNoteTagIndex(ctx, w, site)
	})
}

// HandleNoteTagPage serves a page showing notes for a specific tag.
func (s *Server) HandleNoteTagPage(w http.ResponseWriter, r *http.Request) {
	site := SiteFromContext(r.Context())
	tag := r.PathValue("tag")
	notes := site.Notes.ForTag(tag)

	if notes == nil {
		s.Handle404(w, r)
		return
	}

	s.renderHTML(w, r, "HandleNoteTagPage", func(ctx context.Context, w io.Writer) error {
		return s.Renderer.RenderNoteTagPage(ctx, w, site, tag, notes)
	})
}

// HandleManifest serves the web manifest for the site.
func (s *Server) HandleManifest(w http.ResponseWriter, r *http.Request) {
	site := SiteFromContext(r.Context())
	ctx := r.Context()
	manifest := site.Manifest()

	w.Header().Set("Content-Type", "application/manifest+json")
	if err := s.Renderer.RenderManifest(ctx, w, site, manifest); err != nil {
		log.Printf("ERR: HandleManifest: %s: %v", r.URL.Path, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// HandleRSS serves the RSS feed for the site.
func (s *Server) HandleRSS(w http.ResponseWriter, r *http.Request) {
	site := SiteFromContext(r.Context())
	ctx := r.Context()
	feed := site.RSSFeed()

	w.Header().Set("Content-Type", "application/rss+xml")
	if err := s.Renderer.RenderRSSFeed(ctx, w, site, feed); err != nil {
		log.Printf("ERR: HandleRSS: %s: %v", r.URL.Path, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// HandlePost serves a single post page.
func (s *Server) HandlePost(w http.ResponseWriter, r *http.Request) {
	site := SiteFromContext(r.Context())
	post := site.Posts.GetBySlug(r.PathValue("slug"))
	if post == nil {
		s.Handle404(w, r)
		return
	}

	s.renderHTML(w, r, "HandlePost", func(ctx context.Context, w io.Writer) error {
		return s.Renderer.RenderPost(ctx, w, site, post)
	})
}

// HandleSeriesPost serves a post that belongs to a series.
func (s *Server) HandleSeriesPost(w http.ResponseWriter, r *http.Request) {
	site := SiteFromContext(r.Context())
	seriesSlug := r.PathValue("seriesSlug")
	postSlug := r.PathValue("postSlug")
	fullSlug := seriesSlug + "/" + postSlug

	post := site.Posts.GetBySlug(fullSlug)
	if post == nil {
		s.Handle404(w, r)
		return
	}

	s.renderHTML(w, r, "HandleSeriesPost", func(ctx context.Context, w io.Writer) error {
		return s.Renderer.RenderPost(ctx, w, site, post)
	})
}

// HandleSeriesIndex serves the index page for a series.
func (s *Server) HandleSeriesIndex(w http.ResponseWriter, r *http.Request) {
	site := SiteFromContext(r.Context())
	seriesSlug := r.PathValue("seriesSlug")
	series := site.Series.GetBySlug(seriesSlug)

	if series == nil {
		s.Handle404(w, r)
		return
	}

	s.renderHTML(w, r, "HandleSeriesIndex", func(ctx context.Context, w io.Writer) error {
		return s.Renderer.RenderSeriesIndex(ctx, w, site, series)
	})
}

// HandleTagIndex serves the tags index page.
func (s *Server) HandleTagIndex(w http.ResponseWriter, r *http.Request) {
	site := SiteFromContext(r.Context())
	s.renderHTML(w, r, "HandleTagIndex", func(ctx context.Context, w io.Writer) error {
		return s.Renderer.RenderTagIndex(ctx, w, site)
	})
}

// HandleTagPage serves a page showing posts for a specific tag.
func (s *Server) HandleTagPage(w http.ResponseWriter, r *http.Request) {
	site := SiteFromContext(r.Context())
	tag := r.PathValue("tag")
	posts := site.Posts.ForTag(tag)

	if posts == nil {
		s.Handle404(w, r)
		return
	}

	s.renderHTML(w, r, "HandleTagPage", func(ctx context.Context, w io.Writer) error {
		return s.Renderer.RenderTagPage(ctx, w, site, tag, posts)
	})
}

// HandleArchiveIndex serves the archive index page.
func (s *Server) HandleArchiveIndex(w http.ResponseWriter, r *http.Request) {
	site := SiteFromContext(r.Context())
	s.renderHTML(w, r, "HandleArchiveIndex", func(ctx context.Context, w io.Writer) error {
		return s.Renderer.RenderArchiveIndex(ctx, w, site)
	})
}

// HandleArchivePage serves a page showing posts for a specific year.
func (s *Server) HandleArchivePage(w http.ResponseWriter, r *http.Request) {
	site := SiteFromContext(r.Context())
	year := r.PathValue("year")
	posts := site.Posts.ForYear(year)

	if posts == nil {
		s.Handle404(w, r)
		return
	}

	s.renderHTML(w, r, "HandleArchivePage", func(ctx context.Context, w io.Writer) error {
		return s.Renderer.RenderArchivePage(ctx, w, site, year, posts)
	})
}

// Handle404 serves a custom 404 error page.
func (s *Server) Handle404(w http.ResponseWriter, r *http.Request) {
	site := SiteFromContext(r.Context())
	var buf bytes.Buffer
	if err := s.Renderer.Render404(r.Context(), &buf, site); err != nil {
		log.Printf("ERR: Handle404: %s: %v", r.URL.Path, err)
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusNotFound)
	w.Write(buf.Bytes()) // #nosec G104 - Write errors cannot be handled after headers sent
}
