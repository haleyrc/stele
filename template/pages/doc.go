// Package pages contains templates for individual pages in a rendered blog.
package pages

import (
	"fmt"

	"github.com/a-h/templ"
)

// LayoutFunc represents a function that takes the name of a specific page
// (e.g. "Home") and returns a component that can take other components as
// children; that is: a layout.
type LayoutFunc func(pageName string) templ.Component

func urlf(format string, args ...any) templ.SafeURL {
	s := fmt.Sprintf(format, args...)
	return templ.URL(s)
}
