package site_test

import (
	"strings"
	"testing"

	"github.com/haleyrc/assert"
	"github.com/haleyrc/stele/internal/site"
)

func TestLoadAbout(t *testing.T) {
	about, err := site.LoadAbout("testdata")
	assert.OK(t, err).Fatal()

	if about == nil {
		t.Fatal("expected about to not be nil")
	}

	if !strings.Contains(about.Content, "Some stuff about me") {
		t.Errorf("expected about content to contain 'Some stuff about me', got: %s", about.Content)
	}
}

func TestLoadAbout_Missing(t *testing.T) {
	about, err := site.LoadAbout("/tmp")
	assert.OK(t, err).Fatal()

	if about != nil {
		t.Fatal("expected about to be nil when file doesn't exist")
	}
}

func TestSite_AboutAndSocial(t *testing.T) {
	s, err := site.New("testdata", site.SiteOptions{IncludeDrafts: false})
	assert.OK(t, err).Fatal()

	// Check About is loaded
	if s.About == nil {
		t.Fatal("expected about to not be nil")
	}
	if !strings.Contains(s.About.Content, "Some stuff about me") {
		t.Errorf("expected about content to contain 'Some stuff about me', got: %s", s.About.Content)
	}

	// Check social links are loaded
	assert.Equal(t, "github link", "https://github.com/username", s.Config.Social.GitHub)
	assert.Equal(t, "linkedin link", "https://www.linkedin.com/in/username", s.Config.Social.LinkedIn)
}
