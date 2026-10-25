package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLocalGlob(t *testing.T) {
	dir := t.TempDir()
	font := filepath.Join(dir, "font[1].ttf")
	if err := os.WriteFile(font, []byte("not a font"), 0o600); err != nil {
		t.Fatal(err)
	}

	// Escaped brackets in the pattern must not be expanded again after matching.
	err := InstallFont(filepath.Join(dir, "font[[]1].ttf"))
	if err == nil || !strings.Contains(err.Error(), "cannot parse font") {
		t.Fatalf("expected matched file to be parsed, got %v", err)
	}

	err = InstallFont(filepath.Join(dir, "missing*.ttf"))
	if err == nil || !strings.Contains(err.Error(), "no files match") {
		t.Fatalf("expected unmatched pattern error, got %v", err)
	}
}
