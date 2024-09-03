package stele

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
)

func RenderManifest(ctx context.Context, w io.Writer, cfg *Config) error {
	m, err := newManifest(cfg)
	if err != nil {
		return fmt.Errorf("render manifest: %w", err)
	}

	bytes, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return fmt.Errorf("render manifest: %w", err)
	}

	if _, err := w.Write(bytes); err != nil {
		return fmt.Errorf("render manifest: %w", err)
	}

	return nil
}

type manifest struct {
	BackgroundColor string         `json:"background_color"`
	Categories      []string       `json:"categories"`
	Description     string         `json:"description"`
	Display         string         `json:"display"`
	Icons           []manifestIcon `json:"icons"`
	Name            string         `json:"name"`
	StartURL        string         `json:"start_url"`
}

func newManifest(cfg *Config) (*manifest, error) {
	m := &manifest{
		BackgroundColor: "white",
		Categories:      cfg.Categories,
		Description:     cfg.Description,
		Display:         "fullscreen",
		Icons:           []manifestIcon{},
		Name:            cfg.Title,
		StartURL:        cfg.BaseURL,
	}
	return m, nil
}

type manifestIcon struct {
	Purpose string `json:"purpose"`
	Sizes   string `json:"sizes"`
	Source  string `json:"src"`
	Type    string `json:"type"`
}
