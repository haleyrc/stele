package manifest_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/haleyrc/stele/site"
	"github.com/haleyrc/stele/site/config"
	"github.com/haleyrc/stele/site/manifest"
)

func TestNew(t *testing.T) {
	site := &site.Site{
		Config: &config.Config{
			BaseURL:     "https://example.com",
			Categories:  []string{"go", "react"},
			Description: "This is a test blog",
			Name:        "Test",
		},
	}
	want := manifest.Manifest{
		BackgroundColor: "white",
		Categories:      []string{"go", "react"},
		Description:     "This is a test blog",
		Display:         "fullscreen",
		Icons:           []manifest.Icon{},
		Name:            "Test",
		StartURL:        "https://example.com",
	}

	got := manifest.New(site)
	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("incorrect manifest (-want, +got):\n%s", diff)
	}
}
