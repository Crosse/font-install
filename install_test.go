package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLocalGlob(t *testing.T) {
	t.Parallel()

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

	// An existing literal path must not be treated as a pattern, which would
	// match "font1.ttf" instead. The decoy looks like gzip so reading it fails
	// differently.
	decoy := filepath.Join(dir, "font1.ttf")
	if err := os.WriteFile(decoy, []byte("\x1f\x8b\x08 not gzip"), 0o600); err != nil {
		t.Fatal(err)
	}

	err = InstallFont(font)
	if err == nil || !strings.Contains(err.Error(), "cannot parse font") {
		t.Fatalf("expected literal file to be parsed, got %v", err)
	}

	err = InstallFont(filepath.Join(dir, "missing*.ttf"))
	if err == nil || !strings.Contains(err.Error(), "no files match") {
		t.Fatalf("expected unmatched pattern error, got %v", err)
	}
}
