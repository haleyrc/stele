package components_test

import (
	"testing"
	"time"

	"github.com/haleyrc/stele/internal/testutil"
	"github.com/haleyrc/stele/template/components"
)

func TestPostList(t *testing.T) {
	component := components.PostList(components.PostListProps{
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
	})
	testutil.TestRenderedOutput(t, component, "postlist.html")
}
