package site

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/haleyrc/stele/internal/markdown"
)

// NoteFrontmatter represents the YAML frontmatter in a markdown note file.
type NoteFrontmatter struct {
	// The title of the note.
	Title string `yaml:"title"`

	// A list of tags to associate with the note.
	Tags []string `yaml:"tags"`

	// Whether the note should be pinned to the top of the notes index.
	Pinned bool `yaml:"pinned"`
}

// Validate checks that the frontmatter contains all required fields and that
// field values are valid.
func (fm *NoteFrontmatter) Validate() error {
	if fm.Title == "" {
		return fmt.Errorf("notes must have a title")
	}

	// Tags field is required but can be an empty array
	if fm.Tags == nil {
		return fmt.Errorf("notes must have a tags field")
	}

	return nil
}

// Note represents a living document.
type Note struct {
	// The YAML frontmatter metadata for the note.
	Frontmatter NoteFrontmatter

	// The URL-safe slug for the note. Used to generate note URLs
	// (/notes/{slug}.html).
	Slug string

	// The rendered HTML content of the note.
	Content string
}

// LoadNote loads the file at path and returns the parsed note.
func LoadNote(path string) (*Note, error) {
	var fm NoteFrontmatter
	if err := markdown.ParseFrontmatter(path, &fm); err != nil {
		return nil, fmt.Errorf("load note: %w", err)
	}

	if err := fm.Validate(); err != nil {
		return nil, fmt.Errorf("load note: %s: %w", path, err)
	}

	var content strings.Builder
	if err := markdown.Parse(path, &content); err != nil {
		return nil, fmt.Errorf("load note: %w", err)
	}

	note := &Note{
		Frontmatter: fm,
		Slug:        strings.TrimSuffix(filepath.Base(path), ".md"),
		Content:     content.String(),
	}

	return note, nil
}

// Notes is a slice of Note that implements sort.Interface.
// Notes are sorted alphabetically by title.
type Notes []*Note

// LoadNotes loads all markdown files in the given directory and returns the
// parsed notes. If the directory does not exist, returns an empty slice with
// no error.
func LoadNotes(dir string) (Notes, error) {
	if _, err := os.Stat(dir); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return Notes{}, nil
		}
		return nil, fmt.Errorf("load notes: %w", err)
	}

	paths, err := filepath.Glob(filepath.Join(dir, "*.md"))
	if err != nil {
		return nil, fmt.Errorf("load notes: %w", err)
	}

	var notes Notes
	for _, path := range paths {
		note, err := LoadNote(path)
		if err != nil {
			return nil, fmt.Errorf("load notes: %w", err)
		}
		notes = append(notes, note)
	}

	notes.Sort()
	return notes, nil
}

// Len returns the length of the notes slice.
func (n Notes) Len() int {
	return len(n)
}

// Less reports whether the note at index i should sort before the note at index j.
// Notes are sorted alphabetically by title.
func (n Notes) Less(i, j int) bool {
	return n[i].Frontmatter.Title < n[j].Frontmatter.Title
}

// Swap swaps the notes at indices i and j.
func (n Notes) Swap(i, j int) {
	n[i], n[j] = n[j], n[i]
}

// Sort sorts the notes alphabetically by title.
func (n Notes) Sort() {
	sort.Sort(n)
}

// GetBySlug returns the note with the given slug, or nil if not found.
func (n Notes) GetBySlug(slug string) *Note {
	for i := range n {
		if n[i].Slug == slug {
			return n[i]
		}
	}
	return nil
}

// Pinned returns all notes where Pinned is true, sorted alphabetically by title.
func (n Notes) Pinned() Notes {
	var pinned Notes
	for _, note := range n {
		if note.Frontmatter.Pinned {
			pinned = append(pinned, note)
		}
	}
	return pinned
}

// HasTags returns true if any note has non-empty tags.
func (n Notes) HasTags() bool {
	for _, note := range n {
		if len(note.Frontmatter.Tags) > 0 {
			return true
		}
	}
	return false
}

// IndexByTag returns an index of notes grouped by common tags. A note with
// multiple tags will appear in multiple entries in the index.
func (n Notes) IndexByTag() NoteIndex {
	m := map[string]Notes{}
	for _, note := range n {
		for _, tag := range note.Frontmatter.Tags {
			m[tag] = append(m[tag], note)
		}
	}

	idx := make(NoteIndex, 0, len(m))
	for key, notes := range m {
		idx = append(idx, NoteIndexEntry{
			Key:   key,
			Notes: notes,
		})
	}

	sort.Slice(idx, func(i, j int) bool {
		return idx[i].Key < idx[j].Key
	})

	return idx
}

// ForTag returns all notes that have the given tag. If no notes match, returns
// nil.
func (n Notes) ForTag(tag string) Notes {
	var matches Notes
	for _, note := range n {
		for _, t := range note.Frontmatter.Tags {
			if t == tag {
				matches = append(matches, note)
				break
			}
		}
	}
	if len(matches) == 0 {
		return nil
	}
	return matches
}

// NoteIndex is a slice of entries where each entry contains a set of notes
// that share a common key e.g. tag, etc.
type NoteIndex []NoteIndexEntry

// NoteIndexEntry represents a collection of notes that share a common key e.g.
// tag, etc.
type NoteIndexEntry struct {
	// The shared key for the notes (e.g., tag name).
	Key string

	// The notes that share this key.
	Notes Notes
}
