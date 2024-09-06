package pages_test

import (
	"context"
	"io"

	"github.com/a-h/templ"
)

func EmptyLayout(_ string) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		children := templ.GetChildren(ctx)
		return children.Render(ctx, w)
	})
}
