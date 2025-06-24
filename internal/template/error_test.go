package template_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/haleyrc/assert"
	"github.com/haleyrc/stele/internal/template"
)

func TestTemplateRenderer_Render404(t *testing.T) {
	s := newTestSite()
	renderer := template.NewTemplateRenderer()

	var buf bytes.Buffer
	err := renderer.Render404(context.Background(), &buf, s)
	assert.OK(t, err).Fatal()

	compareGolden(t, buf.String(), "testdata/golden/error_404.html")
}
