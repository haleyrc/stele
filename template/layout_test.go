package template_test

import (
	"testing"

	"github.com/haleyrc/stele/internal/testutil"
	"github.com/haleyrc/stele/template"
)

func TestDefaultLayout(t *testing.T) {
	component := template.DefaultLayout(
		"Test Blog",
		"A test blog.",
		"Ryan Test",
		"2024",
		[]template.MenuLink{
			{Label: "about", Path: "/about"},
		},
	)("Test Page")
	testutil.TestRenderedOutput(t, component, "default-layout.html")
}
