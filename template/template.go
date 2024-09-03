// Package template contains template definitions for the default stele blog
// pages.
package template

import (
	"fmt"

	"github.com/a-h/templ"
)

func urlf(format string, args ...any) templ.SafeURL {
	s := fmt.Sprintf(format, args...)
	return templ.URL(s)
}
