package stele_test

import (
	"context"
	"testing"

	"github.com/haleyrc/stele"
)

func TestBuild(t *testing.T) {
	ctx := context.Background()
	err := stele.Build(ctx, "testdata", "tmp")
	if err != nil {
		t.Fatal(err)
	}
}
