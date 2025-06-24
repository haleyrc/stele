package server

import (
	"sync"

	"github.com/haleyrc/stele/internal/site"
)

// SiteCache provides thread-safe access to a site instance and error state.
// It preserves the last good site when reloads fail.
type SiteCache struct {
	mu        sync.RWMutex
	site      *site.Site
	lastError error
	sourceDir string
	opts      site.SiteOptions
}

// NewSiteCache creates a new cache with an initial site loaded from sourceDir.
func NewSiteCache(sourceDir string, opts site.SiteOptions) (*SiteCache, error) {
	initialSite, err := site.New(sourceDir, opts)
	if err != nil {
		return nil, err
	}

	return &SiteCache{
		site:      initialSite,
		sourceDir: sourceDir,
		opts:      opts,
	}, nil
}

// Get returns the current site and any error from the last update attempt.
func (c *SiteCache) Get() (*site.Site, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.site, c.lastError
}

// Set updates the cached site and error state.
// If err is non-nil, the site is not updated (preserves last good state).
func (c *SiteCache) Set(s *site.Site, err error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if err != nil {
		c.lastError = err
	} else {
		c.site = s
		c.lastError = nil
	}
}

// Reload attempts to reload the site from the source directory.
// Returns an error if reload fails, but preserves the last good site.
func (c *SiteCache) Reload() error {
	newSite, err := site.New(c.sourceDir, c.opts)
	c.Set(newSite, err)
	return err
}
