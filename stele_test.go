package stele_test

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/haleyrc/stele"
)

func TestBuild(t *testing.T) {
	ctx := context.Background()
	srcDir := "testdata"
	dstDir := t.TempDir()

	if err := stele.Build(ctx, srcDir, dstDir); err != nil {
		t.Fatal(err)
	}

	postsPath := filepath.Join(dstDir, "posts", "*.html")
	files, err := filepath.Glob(postsPath)
	if err != nil {
		t.Fatal(err)
	}

	want := 2
	if got := len(files); got != want {
		t.Errorf("expected %d posts, but got %d: %v", want, len(files), files)
	}
}
