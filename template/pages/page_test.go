package pages_test

import (
	"testing"

	"github.com/haleyrc/stele/internal/testutil"
	"github.com/haleyrc/stele/template/pages"
)

func TestPage(t *testing.T) {
	component := pages.Page(EmptyLayout, pages.PageProps{
		Content: "This is a test",
	})
	testutil.TestRenderedOutput(t, component, "page.html")
}
