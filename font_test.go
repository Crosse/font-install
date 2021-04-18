package main

import (
	"fmt"
	"io/ioutil"
	"testing"
)

func TestFontData(t *testing.T) {
	filename := "test_data/open-sans-v18-latin-regular.ttf"
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Fatalf("error opening file")
	}

	fontData, err := NewFontData(filename, data)
	if err != nil {
		t.Errorf("NewFontData shouldn't error")
	}

	if fontData.FileName != filename {
		t.Fatalf("want %q, got %q", filename, fontData.FileName)
	}

	fontName := "Open Sans Regular"
	if fontData.Name != fontName {
		t.Fatalf("want %q, got %q", fontName, fontData.Name)
	}

	fontFamily := "Open Sans"
	if fontData.Family != fontFamily {
		t.Fatalf("want %q, got %q", fontName, fontData.Family)
	}
}

func TestExtensions(t *testing.T) {
	tests := map[string]struct {
		filename string
		valid    bool
		name     string
		family   string
	}{
		"open-sans (eot)":   {filename: "open-sans-v18-latin-regular.eot", valid: false, name: "", family: ""},
		"open-sans (svg)":   {filename: "open-sans-v18-latin-regular.svg", valid: false, name: "", family: ""},
		"open-sans (ttf)":   {filename: "open-sans-v18-latin-regular.ttf", valid: true, name: "Open Sans Regular", family: "Open Sans"},
		"open-sans (woff)":  {filename: "open-sans-v18-latin-regular.woff", valid: false, name: "", family: ""},
		"open-sans (woff2)": {filename: "open-sans-v18-latin-regular.woff2", valid: false, name: "", family: ""},
		"CODE Bold":         {filename: "CODE Bold.otf", valid: true, name: "Code-Bold", family: "Code Bold"},
		"Inconsolata":       {filename: "Inconsolata-Medium.ttf", valid: true, name: "Inconsolata Medium", family: "Inconsolata"},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			filename := fmt.Sprintf("test_data/%s", tc.filename)

			data, err := ioutil.ReadFile(filename)
			if err != nil {
				t.Fatalf("error opening file %v", filename)
			}

			fontData, err := NewFontData(filename, data)
			if tc.valid {
				if err != nil {
					t.Fatalf(`NewFontData(%s, ..) should not return an error`, filename)
				}

				if fontData.FileName != filename {
					t.Fatalf("%s: want filename %q, got %q", name, filename, fontData.FileName)
				}

				if fontData.Name != tc.name {
					t.Fatalf("%s: want font name %q, got %q", name, tc.name, fontData.Name)
				}

				if fontData.Family != tc.family {
					t.Fatalf("%s: want font family %q, got %q", name, tc.family, fontData.Family)
				}
			} else {
				if err == nil {
					t.Fatalf(`NewFontData(%s, ..) should return an error`, filename)
				}
			}
		})
	}
}
