package pages_test

import (
	"testing"
	"time"

	"github.com/haleyrc/stele/internal/testutil"
	"github.com/haleyrc/stele/template/components"
	"github.com/haleyrc/stele/template/pages"
)

func TestIndex(t *testing.T) {
	component := pages.Index(EmptyLayout, pages.IndexProps{
		LatestPost: &components.PostProps{
			Content:   "This is a test",
			Slug:      "test-post",
			Tags:      []string{"go", "react"},
			Timestamp: time.Date(1977, 5, 25, 0, 0, 0, 0, time.UTC),
			Title:     "Test Post",
		},
		RecentPosts: components.PostListProps{
			Posts: []components.PostListEntryProps{
				{
					Slug:      "test-post-1",
					Timestamp: time.Date(1977, 5, 25, 0, 0, 0, 0, time.UTC),
					Title:     "Test Post 1",
				},
				{
					Slug:      "test-post-2",
					Timestamp: time.Date(1977, 5, 26, 0, 0, 0, 0, time.UTC),
					Title:     "Test Post 2",
				},
			},
		},
	})
	testutil.TestRenderedOutput(t, component, "index.html")
}
