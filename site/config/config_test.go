package config_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/haleyrc/stele/site/config"
)

func TestNew(t *testing.T) {
	want := &config.Config{
		Author:      "Grace Hopper",
		BaseURL:     "https://example.com",
		Categories:  []string{"programming", "music"},
		Description: "This is a test blog",
		Menu: []config.Link{
			{Label: "About", Path: "/about"},
			{Label: "Contact Us", Path: "/contact-us"},
		},
		Name: "Test",
	}

	got, err := config.New("testdata/config.yml")
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("incorrect index (-want, +got):\n%s", diff)
	}
}
