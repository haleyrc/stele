package manifest

import "github.com/haleyrc/stele/site"

type Icon struct {
	Purpose string `json:"purpose"`
	Sizes   string `json:"sizes"`
	Source  string `json:"src"`
	Type    string `json:"type"`
}

type Manifest struct {
	BackgroundColor string   `json:"background_color"`
	Categories      []string `json:"categories"`
	Description     string   `json:"description"`
	Display         string   `json:"display"`
	Icons           []Icon   `json:"icons"`
	Name            string   `json:"name"`
	StartURL        string   `json:"start_url"`
}

func New(s *site.Site) Manifest {
	return Manifest{
		BackgroundColor: "white",
		Categories:      s.Config.Categories,
		Description:     s.Config.Description,
		Display:         "fullscreen",
		Icons:           []Icon{},
		Name:            s.Config.Name,
		StartURL:        s.Config.BaseURL,
	}
}
