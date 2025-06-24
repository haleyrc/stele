package server

import (
	"context"

	"github.com/haleyrc/stele/internal/site"
)

type contextKey int

const siteKey contextKey = 0

// WithSite injects a site into the request context.
func WithSite(ctx context.Context, s *site.Site) context.Context {
	return context.WithValue(ctx, siteKey, s)
}

// SiteFromContext retrieves the site from the request context.
// Panics if site is not in context (indicates middleware failure).
func SiteFromContext(ctx context.Context) *site.Site {
	s, ok := ctx.Value(siteKey).(*site.Site)
	if !ok {
		panic("server: site not found in context")
	}
	return s
}
