package testutil

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/a-h/templ"
	"github.com/a-h/templ/generator/htmldiff"
)

func TestRenderedOutput(t *testing.T, c templ.Component, filename string) {
	want := readGoldenFile(t, filename)

	diff, err := htmldiff.Diff(c, want)
	if err != nil {
		t.Fatal(err)
	}
	if diff != "" {
		t.Error("invalid output (-want +got):\n", diff)
	}
}

func readGoldenFile(t *testing.T, filename string) string {
	path := filepath.Join("testdata", filename)
	bytes, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(bytes)
}
