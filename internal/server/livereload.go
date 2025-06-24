package server

import (
	"context"
	"fmt"
	"html"
	"log"
	"net/http"
	"sync"
	"time"

	httpserver "github.com/haleyrc/server"
	"github.com/haleyrc/stele/internal/site"
)

// LiveReloader wraps an http.Handler to add live reload capabilities.
// It intercepts dev endpoints, manages site reloading, and injects
// site into context for the wrapped handler.
type LiveReloader struct {
	port    string
	handler http.Handler
	cache   *SiteCache
	watcher *Watcher

	// SSE fields (inlined)
	mu      sync.RWMutex
	clients map[chan string]struct{}
	ctx     context.Context
	cancel  context.CancelFunc
}

// NewLiveReloader creates a new live reload wrapper.
func NewLiveReloader(
	port string,
	renderer site.Renderer,
	cache *SiteCache,
) (*LiveReloader, error) {
	server := NewServer(renderer)

	lr := &LiveReloader{
		port:    port,
		handler: server,
		cache:   cache,
		clients: make(map[chan string]struct{}),
	}

	watcher, err := NewWatcher(cache.sourceDir, lr.reload)
	if err != nil {
		return nil, err
	}
	lr.watcher = watcher

	return lr, nil
}

// Start begins watching for file changes.
func (lr *LiveReloader) Start(ctx context.Context) {
	lr.ctx, lr.cancel = context.WithCancel(ctx)
	lr.watcher.Start(lr.ctx)

	// Close SSE connections on context cancellation
	go func() {
		<-lr.ctx.Done()
		lr.closeAll()
	}()
}

// ListenAndServe starts the HTTP server with live reload enabled.
func (lr *LiveReloader) ListenAndServe(ctx context.Context) error {
	log.Printf("Listening on http://localhost:%s", lr.port)
	srv := httpserver.New(lr.port, lr)
	return srv.ListenAndServe(ctx)
}

// ServeHTTP implements http.Handler.
func (lr *LiveReloader) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Intercept dev endpoints
	switch r.URL.Path {
	case "/__dev__/sse":
		lr.handleSSE(w, r)
		return
	case "/__dev__/reload.js":
		lr.handleScript(w, r)
		return
	}

	// Get site and check for errors
	site, err := lr.cache.Get()
	if err != nil {
		lr.renderErrorPage(w, r, err)
		return
	}

	// Inject site into context and delegate to wrapped handler
	ctx := WithSite(r.Context(), site)
	lr.handler.ServeHTTP(w, r.WithContext(ctx))
}

// reload attempts to reload the site and broadcasts the result.
func (lr *LiveReloader) reload() {
	log.Println("Reloading site content...")

	err := lr.cache.Reload()

	if err != nil {
		log.Printf("ERR: reload failed: %v", err)
	} else {
		log.Println("Site reloaded successfully")
	}

	// Broadcast reload event (even on error, so error page displays)
	lr.broadcast("reload")
}

// renderErrorPage renders a simple error page when site fails to reload.
func (lr *LiveReloader) renderErrorPage(w http.ResponseWriter, r *http.Request, err error) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusInternalServerError)

	// Simple inline template (no need for fancy rendering)
	fmt.Fprintf(w, `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Reload Error</title>
    <script src="/__dev__/reload.js"></script>
    <style>
        body { font-family: monospace; padding: 2rem; max-width: 800px; margin: 0 auto; }
        h1 { color: #c00; }
        pre { background: #f4f4f4; padding: 1rem; overflow-x: auto; }
    </style>
</head>
<body>
    <h1>Site failed to reload</h1>
    <p>The site encountered an error during reload. The previous version is still being served.</p>
    <p>Fix the error and save to reload automatically.</p>
    <pre>%s</pre>
</body>
</html>`, html.EscapeString(err.Error()))
}

// handleScript serves the live reload JavaScript.
func (lr *LiveReloader) handleScript(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/javascript")
	fmt.Fprint(w, `(function() {
    const source = new EventSource('/__dev__/sse');

    source.onmessage = function(event) {
        if (event.data === 'reload') {
            console.log('[stele] Reloading...');
            window.location.reload();
        }
    };

    source.onerror = function() {
        console.log('[stele] SSE connection lost, retrying...');
    };

    console.log('[stele] Live reload connected');
})();`)
}

// SSE implementation (inlined from SSEBroadcaster)

func (lr *LiveReloader) handleSSE(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	rc := http.NewResponseController(w)
	_ = rc.SetReadDeadline(time.Time{})  // #nosec G104 - Optional deadline setting
	_ = rc.SetWriteDeadline(time.Time{}) // #nosec G104 - Optional deadline setting

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "SSE not supported", http.StatusInternalServerError)
		return
	}

	events := lr.subscribe()
	defer lr.unsubscribe(events)

	for {
		select {
		case <-r.Context().Done():
			return
		case <-lr.ctx.Done():
			return
		case event := <-events:
			fmt.Fprintf(w, "data: %s\n\n", event)
			flusher.Flush()
		}
	}
}

func (lr *LiveReloader) subscribe() <-chan string {
	ch := make(chan string, 1)

	lr.mu.Lock()
	lr.clients[ch] = struct{}{}
	lr.mu.Unlock()

	return ch
}

func (lr *LiveReloader) unsubscribe(ch <-chan string) {
	lr.mu.Lock()
	writeCh, ok := (interface{})(ch).(chan string)
	if ok {
		delete(lr.clients, writeCh)
		close(writeCh)
	}
	lr.mu.Unlock()
}

func (lr *LiveReloader) broadcast(event string) {
	lr.mu.RLock()
	defer lr.mu.RUnlock()

	for ch := range lr.clients {
		select {
		case ch <- event:
		default:
			// Client too slow, skip
		}
	}
}

func (lr *LiveReloader) closeAll() {
	lr.mu.Lock()
	defer lr.mu.Unlock()

	for ch := range lr.clients {
		close(ch)
	}
	lr.clients = make(map[chan string]struct{})
}
