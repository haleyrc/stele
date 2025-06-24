package site_test

import (
	"testing"

	"github.com/haleyrc/assert"
	"github.com/haleyrc/stele/internal/site"
)

func TestLoadNote(t *testing.T) {
	note, err := site.LoadNote("testdata/notes/golang-tips.md")
	assert.OK(t, err).Fatal()

	// Verify frontmatter fields are parsed correctly
	assert.Equal(t, "title", "Golang Best Practices", note.Frontmatter.Title)
	assert.Equal(t, "pinned", true, note.Frontmatter.Pinned)
	assert.Equal(t, "slug", "golang-tips", note.Slug)

	expectedTags := []string{"golang", "programming"}
	assert.SliceEqual(t, "tags", expectedTags, note.Frontmatter.Tags)
}

func TestLoadNote_NotPinned(t *testing.T) {
	note, err := site.LoadNote("testdata/notes/vim-shortcuts.md")
	assert.OK(t, err).Fatal()

	assert.Equal(t, "title", "Vim Shortcuts", note.Frontmatter.Title)
	assert.Equal(t, "pinned", false, note.Frontmatter.Pinned)
}

func TestNotesSorting(t *testing.T) {
	// Load notes from testdata
	notes := site.Notes{}

	noteFiles := []string{
		"testdata/notes/golang-tips.md",
		"testdata/notes/vim-shortcuts.md",
		"testdata/notes/algorithms.md",
	}

	for _, file := range noteFiles {
		note, err := site.LoadNote(file)
		assert.OK(t, err).Fatal()
		notes = append(notes, note)
	}

	// Sort notes
	notes.Sort()

	// Verify notes are sorted alphabetically by title
	assert.Equal(t, "first note title", "Algorithm Notes", notes[0].Frontmatter.Title)
	assert.Equal(t, "second note title", "Golang Best Practices", notes[1].Frontmatter.Title)
	assert.Equal(t, "third note title", "Vim Shortcuts", notes[2].Frontmatter.Title)

	// Verify titles are indeed in alphabetical order
	for i := 0; i < len(notes)-1; i++ {
		if notes[i].Frontmatter.Title > notes[i+1].Frontmatter.Title {
			t.Errorf("notes[%d] title should be before notes[%d] title", i, i+1)
		}
	}
}

func TestNotes_Pinned(t *testing.T) {
	// Load notes from testdata
	notes := site.Notes{}

	noteFiles := []string{
		"testdata/notes/golang-tips.md",   // pinned: true
		"testdata/notes/vim-shortcuts.md", // pinned: false
		"testdata/notes/algorithms.md",    // pinned: false (default)
	}

	for _, file := range noteFiles {
		note, err := site.LoadNote(file)
		assert.OK(t, err).Fatal()
		notes = append(notes, note)
	}

	pinned := notes.Pinned()

	assert.Equal(t, "pinned count", 1, len(pinned))
	assert.Equal(t, "pinned note title", "Golang Best Practices", pinned[0].Frontmatter.Title)
}

func TestNotes_GetBySlug(t *testing.T) {
	notes := site.Notes{}

	noteFiles := []string{
		"testdata/notes/golang-tips.md",
		"testdata/notes/vim-shortcuts.md",
	}

	for _, file := range noteFiles {
		note, err := site.LoadNote(file)
		assert.OK(t, err).Fatal()
		notes = append(notes, note)
	}

	found := notes.GetBySlug("golang-tips")
	assert.Equal(t, "found note title", "Golang Best Practices", found.Frontmatter.Title)

	notFound := notes.GetBySlug("nonexistent")
	if notFound != nil {
		t.Error("expected nil for nonexistent slug")
	}
}

func TestNotes_IndexByTag(t *testing.T) {
	notes := site.Notes{}

	noteFiles := []string{
		"testdata/notes/golang-tips.md",   // tags: golang, programming
		"testdata/notes/vim-shortcuts.md", // tags: vim, productivity
		"testdata/notes/algorithms.md",    // tags: algorithms, programming
	}

	for _, file := range noteFiles {
		note, err := site.LoadNote(file)
		assert.OK(t, err).Fatal()
		notes = append(notes, note)
	}

	index := notes.IndexByTag()

	// Verify index entries are sorted by key (tag name)
	expectedTags := []string{"algorithms", "golang", "productivity", "programming", "vim"}
	for i, entry := range index {
		assert.Equal(t, "tag key", expectedTags[i], entry.Key)
	}

	// Verify programming tag has 2 notes
	var programmingEntry *site.NoteIndexEntry
	for _, entry := range index {
		if entry.Key == "programming" {
			programmingEntry = &entry
			break
		}
	}

	if programmingEntry == nil {
		t.Fatal("expected to find programming tag in index")
	}

	assert.Equal(t, "programming tag count", 2, len(programmingEntry.Notes))
}

func TestNotes_ForTag(t *testing.T) {
	notes := site.Notes{}

	noteFiles := []string{
		"testdata/notes/golang-tips.md",
		"testdata/notes/vim-shortcuts.md",
		"testdata/notes/algorithms.md",
	}

	for _, file := range noteFiles {
		note, err := site.LoadNote(file)
		assert.OK(t, err).Fatal()
		notes = append(notes, note)
	}

	programmingNotes := notes.ForTag("programming")
	assert.Equal(t, "programming notes count", 2, len(programmingNotes))

	vimNotes := notes.ForTag("vim")
	assert.Equal(t, "vim notes count", 1, len(vimNotes))

	nonexistentNotes := notes.ForTag("nonexistent")
	if nonexistentNotes != nil {
		t.Error("expected nil for nonexistent tag")
	}
}

func TestNotes_HasTags(t *testing.T) {
	notes := site.Notes{}

	noteFiles := []string{
		"testdata/notes/golang-tips.md",
	}

	for _, file := range noteFiles {
		note, err := site.LoadNote(file)
		assert.OK(t, err).Fatal()
		notes = append(notes, note)
	}

	if !notes.HasTags() {
		t.Error("expected notes to have tags")
	}

	emptyNotes := site.Notes{}
	if emptyNotes.HasTags() {
		t.Error("expected empty notes to not have tags")
	}
}

func TestLoadNotes(t *testing.T) {
	notes, err := site.LoadNotes("testdata/notes")
	assert.OK(t, err).Fatal()

	assert.Equal(t, "notes count", 3, len(notes))

	// Verify notes are sorted alphabetically by title
	assert.Equal(t, "first note title", "Algorithm Notes", notes[0].Frontmatter.Title)
	assert.Equal(t, "second note title", "Golang Best Practices", notes[1].Frontmatter.Title)
	assert.Equal(t, "third note title", "Vim Shortcuts", notes[2].Frontmatter.Title)
}

func TestLoadNotes_NonexistentDirectory(t *testing.T) {
	notes, err := site.LoadNotes("testdata/nonexistent")
	assert.OK(t, err).Fatal()

	assert.Equal(t, "notes count", 0, len(notes))
}
