package site_test

import (
	"testing"

	"github.com/haleyrc/assert"
	"github.com/haleyrc/stele/internal/site"
)

// assertConfigMatchesTestdata verifies that the config matches the expected values from testdata/site.yaml
func assertConfigMatchesTestdata(t *testing.T, config site.SiteConfig) {
	assert.Equal(t, "author", "Ryan Haley", config.Author)
	assert.Equal(t, "base URL", "https://blog.ryanchaley.com", config.BaseURL)
	assert.Equal(t, "description", "A barely coherent train of thoughts", config.Description)
	assert.Equal(t, "title", "Taints and Tolerations", config.Title)

	// Verify categories array
	expectedCategories := []string{"blog", "programming", "dev", "development"}
	assert.SliceEqual(t, "categories", expectedCategories, config.Categories)

	// Verify social links
	assert.Equal(t, "github", "https://github.com/username", config.Social.GitHub)
	assert.Equal(t, "linkedin", "https://www.linkedin.com/in/username", config.Social.LinkedIn)
}

func assertAllPostsLoaded(t *testing.T, site *site.Site) {
	// Posts should be sorted by timestamp in descending order (newest first)
	// Includes both standalone posts and series posts
	expectedSlugs := []string{
		"draft-exploring-go-generics",  // 2025-10-01 (newest, draft)
		"getting-started-with-go",      // 2025-09-20
		"building-rest-apis-go",        // 2025-09-18
		"advanced-go-patterns",         // 2025-09-15
		"testing-in-go-complete-guide", // 2025-09-12
		"go-basics/functions",          // 2024-02-01 (from series)
		"go-basics/variables",          // 2024-01-01 (from series, oldest)
	}

	actualSlugs := make([]string, len(site.Posts))
	for i, post := range site.Posts {
		actualSlugs[i] = post.Slug
	}

	assert.SliceEqual(t, "post slugs", expectedSlugs, actualSlugs)
}

func TestNewSite(t *testing.T) {
	site, err := site.New("testdata", site.SiteOptions{IncludeDrafts: true})
	assert.OK(t, err).Fatal()

	// Verify the site directory is set correctly
	assert.Equal(t, "site directory", "testdata", site.Dir)

	assertConfigMatchesTestdata(t, site.Config)
	assertAllPostsLoaded(t, site)
}

func TestLoadSiteConfig(t *testing.T) {
	config, err := site.LoadSiteConfig("testdata")
	assert.OK(t, err).Fatal()

	assertConfigMatchesTestdata(t, *config)
}

func TestSite_LatestPost(t *testing.T) {
	site, err := site.New("testdata", site.SiteOptions{IncludeDrafts: true})
	assert.OK(t, err).Fatal()

	latestPost := site.Posts.Latest()
	if latestPost == nil {
		t.Fatal("expected latest post to not be nil")
	}
	assert.Equal(t, "latest post slug", "draft-exploring-go-generics", latestPost.Slug)
}

func TestSite_LatestPost_EmptySite(t *testing.T) {
	s := &site.Site{Posts: site.Posts{}}

	latestPost := s.Posts.Latest()
	if latestPost != nil {
		t.Error("expected latest post from empty site to be nil")
	}
}

func TestSite_RecentPosts(t *testing.T) {
	site, err := site.New("testdata", site.SiteOptions{IncludeDrafts: true})
	assert.OK(t, err).Fatal()

	// Test getting 2 recent posts (should include the latest)
	recentPosts := site.Posts.Recent(2)
	assert.Equal(t, "recent posts count", 2, len(recentPosts))
	assert.Equal(t, "first recent post", "draft-exploring-go-generics", recentPosts[0].Slug)
	assert.Equal(t, "second recent post", "getting-started-with-go", recentPosts[1].Slug)
}

func TestSite_RecentPosts_AllPosts(t *testing.T) {
	site, err := site.New("testdata", site.SiteOptions{IncludeDrafts: true})
	assert.OK(t, err).Fatal()

	// Test getting all recent posts (maxCount = 0)
	// Now includes series posts as well
	recentPosts := site.Posts.Recent(0)
	assert.Equal(t, "recent posts count", 7, len(recentPosts))
	assert.Equal(t, "first recent post", "draft-exploring-go-generics", recentPosts[0].Slug)
	assert.Equal(t, "second recent post", "getting-started-with-go", recentPosts[1].Slug)
	assert.Equal(t, "third recent post", "building-rest-apis-go", recentPosts[2].Slug)
	assert.Equal(t, "fourth recent post", "advanced-go-patterns", recentPosts[3].Slug)
	assert.Equal(t, "fifth recent post", "testing-in-go-complete-guide", recentPosts[4].Slug)
}

func TestSite_RecentPosts_EmptySite(t *testing.T) {
	s := &site.Site{Posts: site.Posts{}}

	recentPosts := s.Posts.Recent(5)
	if recentPosts != nil {
		t.Error("expected recent posts from empty site to be nil")
	}
}

func TestSite_RecentPosts_SinglePost(t *testing.T) {
	s := &site.Site{Posts: site.Posts{{}}}

	recentPosts := s.Posts.Recent(5)
	assert.Equal(t, "single post site count", 1, len(recentPosts))
}
