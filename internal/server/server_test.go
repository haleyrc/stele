package server_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/haleyrc/assert"
	"github.com/haleyrc/stele/internal/server"
	"github.com/haleyrc/stele/internal/template"
	"github.com/haleyrc/stele/internal/testutil"
)

func TestServer_HandleManifest(t *testing.T) {
	s := testutil.TestSite()
	renderer := template.NewTemplateRenderer()
	srv := server.NewServer(renderer)

	req := httptest.NewRequest("GET", "/manifest.webmanifest", nil)
	req = req.WithContext(server.WithSite(req.Context(), s))
	rr := httptest.NewRecorder()

	srv.HandleManifest(rr, req)

	assert.Equal(t, "status code", http.StatusOK, rr.Code)
	assert.Equal(t, "content type", "application/manifest+json", rr.Header().Get("Content-Type"))

	testutil.AssertRenderedOutput(t, testutil.ExpectedManifest, rr.Body.String())
}

func TestServer_HandleRSS(t *testing.T) {
	s := testutil.TestSite()
	renderer := template.NewTemplateRenderer()
	srv := server.NewServer(renderer)

	req := httptest.NewRequest("GET", "/rss.xml", nil)
	req = req.WithContext(server.WithSite(req.Context(), s))
	rr := httptest.NewRecorder()

	srv.HandleRSS(rr, req)

	assert.Equal(t, "status code", http.StatusOK, rr.Code)
	assert.Equal(t, "content type", "application/rss+xml", rr.Header().Get("Content-Type"))

	// Verify RSS XML contains expected content
	body := rr.Body.String()
	if !strings.Contains(body, `<?xml version="1.0" encoding="UTF-8" ?>`) {
		t.Error("expected RSS to contain XML declaration")
	}
	if !strings.Contains(body, "<rss version=\"2.0\"") {
		t.Error("expected RSS to contain rss element with version")
	}
	if !strings.Contains(body, "<title>Alice Codes</title>") {
		t.Error("expected RSS to contain site title")
	}
	if !strings.Contains(body, "<link>https://alice.dev</link>") {
		t.Error("expected RSS to contain site link")
	}
}

func TestServer_HandlePost(t *testing.T) {
	s := testutil.TestSite()
	renderer := template.NewTemplateRenderer()
	srv := server.NewServer(renderer)

	t.Run("existing post", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/posts/getting-started-with-go.html", nil)
		req = req.WithContext(server.WithSite(req.Context(), s))
		req.SetPathValue("slug", "getting-started-with-go")
		rr := httptest.NewRecorder()

		srv.HandlePost(rr, req)

		assert.Equal(t, "status code", http.StatusOK, rr.Code)
		assert.Equal(t, "content type", "text/html; charset=utf-8", rr.Header().Get("Content-Type"))

		body := rr.Body.String()
		if !strings.Contains(body, "Getting Started with Go") {
			t.Errorf("expected page to contain post title 'Getting Started with Go', but it didn't")
		}
	})

	t.Run("non-existent post", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/posts/non-existent-post.html", nil)
		req = req.WithContext(server.WithSite(req.Context(), s))
		req.SetPathValue("slug", "non-existent-post")
		rr := httptest.NewRecorder()

		srv.HandlePost(rr, req)

		assert.Equal(t, "status code", http.StatusNotFound, rr.Code)
	})
}
