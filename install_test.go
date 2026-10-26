package main

import (
	"net/url"
	"os"
	"path/filepath"
	"runtime"
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

// TestPathForms checks that each kind of argument reaches the file it names.
// The files aren't fonts, so reaching them fails with "cannot parse font".
func TestPathForms(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()

	// Names that url.Parse would misread as a scheme or an escape.
	names := []string{"100%zz.ttf"}
	if runtime.GOOS != "windows" {
		names = append(names, "font:bold.ttf")
	}

	for _, name := range names {
		font := filepath.Join(dir, name)
		if err := os.WriteFile(font, []byte("not a font"), 0o600); err != nil {
			t.Fatal(err)
		}

		err := InstallFont(font)
		if err == nil || !strings.Contains(err.Error(), "cannot parse font") {
			t.Errorf("%s: expected file to be parsed, got %v", name, err)
		}
	}

	font := filepath.Join(dir, "plain.ttf")
	if err := os.WriteFile(font, []byte("not a font"), 0o600); err != nil {
		t.Fatal(err)
	}

	urlPath := filepath.ToSlash(font)
	if !strings.HasPrefix(urlPath, "/") {
		urlPath = "/" + urlPath
	}

	fileURL := (&url.URL{Scheme: "file", Path: urlPath}).String()

	err := InstallFont(fileURL)
	if err == nil || !strings.Contains(err.Error(), "cannot parse font") {
		t.Errorf("%s: expected file to be parsed, got %v", fileURL, err)
	}

	err = InstallFont("ftp://example.com/font.ttf")
	if err == nil || !strings.Contains(err.Error(), "unhandled URL scheme") {
		t.Errorf("expected unhandled scheme error, got %v", err)
	}
}

func TestFileURLToPath(t *testing.T) {
	t.Parallel()

	const (
		unixFont    = "/usr/share/fonts/x.ttf"
		windowsFont = `C:\fonts\x.ttf`
	)

	tests := map[string]string{
		"file:///usr/share/fonts/x.ttf":          unixFont,
		"file://localhost/usr/share/fonts/x.ttf": unixFont,
		"file://LOCALHOST/usr/share/fonts/x.ttf": unixFont,
		"file:///tmp/My%20Fonts/x.ttf":           "/tmp/My Fonts/x.ttf",
		"file://server/share/x.ttf":              "",
	}
	if runtime.GOOS == "windows" {
		tests = map[string]string{
			"file:///C:/fonts/x.ttf":          windowsFont,
			"file://localhost/C:/fonts/x.ttf": windowsFont,
			"file://LOCALHOST/C:/fonts/x.ttf": windowsFont,
			"file:///C:/My%20Fonts/x.ttf":     `C:\My Fonts\x.ttf`,
			"file://server/share/x.ttf":       `\\server\share\x.ttf`,
		}
	}

	for rawURL, want := range tests {
		u, err := url.Parse(rawURL)
		if err != nil {
			t.Fatal(err)
		}

		got, err := fileURLToPath(u)

		switch {
		case want == "" && err == nil:
			t.Errorf("%s: expected an error, got %q", rawURL, got)
		case want != "" && (err != nil || got != want):
			t.Errorf("%s: got %q, %v; want %q", rawURL, got, err, want)
		}
	}
}
