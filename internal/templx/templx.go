// Package templx provides convenience functions for working with templ
// templates.
package templx

import (
	"fmt"

	"github.com/a-h/templ"
)

// URLf formats a URL string and returns it as a templ.SafeURL, indicating
// that the content is safe to use in URL contexts without escaping.
func URLf(format string, args ...any) templ.SafeURL {
	s := fmt.Sprintf(format, args...)
	return templ.URL(s)
}
