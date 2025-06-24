package site

import (
	"encoding/json"
	"fmt"
	"io"
)

// Manifest represents a Progressive Web App manifest for a site.
type Manifest struct {
	// The default background color for the application.
	BackgroundColor string `json:"background_color"`

	// One or more categories that the application belongs to.
	Categories []string `json:"categories"`

	// A brief description of the application.
	Description string `json:"description"`

	// The preferred display mode for the application (e.g., "fullscreen",
	// "standalone", "minimal-ui", "browser").
	Display string `json:"display"`

	// An array of image objects that can serve as application icons.
	Icons []ManifestIcon `json:"icons"`

	// The name of the application.
	Name string `json:"name"`

	// The URL that loads when a user launches the application.
	StartURL string `json:"start_url"`
}

// NewManifest creates a new web manifest for the given site.
func NewManifest(s *Site) *Manifest {
	manifest := &Manifest{
		BackgroundColor: "white",
		Categories:      s.Config.Categories,
		Description:     s.Config.Description,
		Display:         "fullscreen",
		Icons:           []ManifestIcon{},
		Name:            s.Config.Title,
		StartURL:        s.Config.BaseURL,
	}
	return manifest
}

// Render writes the manifest as JSON to the provided writer.
func (m *Manifest) Render(w io.Writer) error {
	bytes, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return fmt.Errorf("manifest: render: %w", err)
	}

	if _, err := w.Write(bytes); err != nil {
		return fmt.Errorf("manifest: render: %w", err)
	}

	return nil
}

// ManifestIcon represents an application icon in a web manifest.
type ManifestIcon struct {
	// The purpose of the icon (e.g., "maskable", "any").
	Purpose string `json:"purpose"`

	// The sizes of the icon (e.g., "192x192", "512x512").
	Sizes string `json:"sizes"`

	// The URL of the icon image.
	Source string `json:"src"`

	// The media type of the icon image (e.g., "image/png").
	Type string `json:"type"`
}
