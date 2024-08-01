package site

import (
	"fmt"
	"path/filepath"

	"github.com/haleyrc/stele/config"
	"github.com/haleyrc/stele/index"
)

type Site struct {
	Config *config.Config
	Index  *index.Index
}

func New(dir string) (*Site, error) {
	cfg, err := config.New(filepath.Join(dir, "config.yml"))
	if err != nil {
		return nil, errf("site: new", err)
	}

	idx, err := index.New(dir)
	if err != nil {
		return nil, errf("site: new", err)
	}

	site := &Site{
		Config: cfg,
		Index:  idx,
	}

	return site, nil
}

func errf(op string, err error) error {
	return fmt.Errorf("%s: %w", op, err)
}
