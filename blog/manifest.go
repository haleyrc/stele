package blog

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
)

type Manifest struct {
	BackgroundColor string         `json:"background_color"`
	Categories      []string       `json:"categories"`
	Description     string         `json:"description"`
	Display         string         `json:"display"`
	Icons           []ManifestIcon `json:"icons"`
	Name            string         `json:"name"`
	StartURL        string         `json:"start_url"`
}

func NewManifest(cfg *Config) (*Manifest, error) {
	m := &Manifest{
		BackgroundColor: "white",
		Categories:      cfg.Categories,
		Description:     cfg.Description,
		Display:         "fullscreen",
		Icons:           []ManifestIcon{},
		Name:            cfg.Title,
		StartURL:        cfg.BaseURL,
	}
	return m, nil
}

func (m *Manifest) Render(ctx context.Context, w io.Writer) error {
	bytes, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return fmt.Errorf("manifest: render: %w", err)
	}

	if _, err := w.Write(bytes); err != nil {
		return fmt.Errorf("manifest: render: %w", err)
	}

	return nil
}

type ManifestIcon struct {
	Purpose string `json:"purpose"`
	Sizes   string `json:"sizes"`
	Source  string `json:"src"`
	Type    string `json:"type"`
}
