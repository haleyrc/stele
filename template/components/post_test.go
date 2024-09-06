package components_test

import (
	"testing"
	"time"

	"github.com/haleyrc/stele/internal/testutil"
	"github.com/haleyrc/stele/template/components"
)

func TestPost(t *testing.T) {
	component := components.Post(components.PostProps{
		Content:   "This is a test",
		Slug:      "test-post",
		Tags:      []string{"go", "react"},
		Timestamp: time.Date(1977, 5, 25, 0, 0, 0, 0, time.UTC),
		Title:     "Test Post",
	})
	testutil.TestRenderedOutput(t, component, "post.html")
}
