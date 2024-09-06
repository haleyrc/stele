package pages_test

import (
	"testing"

	"github.com/haleyrc/stele/internal/testutil"
	"github.com/haleyrc/stele/template/pages"
)

func TestPostIndex(t *testing.T) {
	component := pages.PostIndex(EmptyLayout, pages.PostIndexProps{
		Entries: []pages.PostIndexEntryProps{
			{Count: 2, Key: "go"},
			{Count: 3, Key: "react"},
		},
		Prefix: "/tags/",
	})
	testutil.TestRenderedOutput(t, component, "postindex.html")
}
