package server

import (
	"context"
	"log"
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"
)

// Watcher monitors filesystem changes and triggers reload callbacks.
type Watcher struct {
	watcher  *fsnotify.Watcher
	onChange func()
}

// NewWatcher creates a new file watcher for the given site directory.
// It watches for changes to content files (posts, notes, config) and invokes
// the onChange callback when relevant files are modified.
func NewWatcher(siteDir string, onChange func()) (*Watcher, error) {
	fw, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	// Watch directories (fsnotify watches recursively on some platforms)
	dirs := []string{
		filepath.Join(siteDir, "posts"),
		filepath.Join(siteDir, "notes"),
		siteDir, // For about.md and stele.yaml
	}

	for _, dir := range dirs {
		if err := fw.Add(dir); err != nil {
			_ = fw.Close() // #nosec G104 - Cleanup error not actionable
			return nil, err
		}
	}

	return &Watcher{
		watcher:  fw,
		onChange: onChange,
	}, nil
}

// Start begins watching for file changes. It runs in a goroutine and
// continues until the context is cancelled.
func (w *Watcher) Start(ctx context.Context) {
	go func() {
		debounce := time.NewTimer(100 * time.Millisecond)
		debounce.Stop()

		defer func() {
			if !debounce.Stop() {
				// Drain the channel if stop returned false
				select {
				case <-debounce.C:
				default:
				}
			}
		}()

		for {
			select {
			case <-ctx.Done():
				_ = w.watcher.Close() // #nosec G104 - Cleanup error during shutdown
				return

			case event, ok := <-w.watcher.Events:
				if !ok {
					return
				}

				// Filter to relevant file types
				if w.isRelevantFile(event.Name) {
					// Debounce rapid-fire saves
					debounce.Reset(100 * time.Millisecond)
				}

			case err, ok := <-w.watcher.Errors:
				if !ok {
					return
				}
				log.Printf("watcher error: %v", err)

			case <-debounce.C:
				w.onChange()
			}
		}
	}()
}

// isRelevantFile checks if the given file path should trigger a reload.
func (w *Watcher) isRelevantFile(path string) bool {
	ext := filepath.Ext(path)
	base := filepath.Base(path)

	// .md files, index.yaml, stele.yaml, about.md
	return ext == ".md" ||
		base == "index.yaml" ||
		base == "stele.yaml"
}
