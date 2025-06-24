package site_test

import (
	"bytes"
	"testing"

	"github.com/haleyrc/assert"
	"github.com/haleyrc/stele/internal/site"
	"github.com/haleyrc/stele/internal/testutil"
)

func TestManifest_Render(t *testing.T) {
	s := testutil.TestSite()
	manifest := site.NewManifest(s)

	var buff bytes.Buffer
	err := manifest.Render(&buff)
	assert.OK(t, err).Fatal()

	testutil.AssertRenderedOutput(t, testutil.ExpectedManifest, buff.String())
}

func TestSite_RenderManifest(t *testing.T) {
	s := testutil.TestSite()

	var buff bytes.Buffer
	manifest := s.Manifest()
	err := manifest.Render(&buff)
	assert.OK(t, err).Fatal()

	testutil.AssertRenderedOutput(t, testutil.ExpectedManifest, buff.String())
}
