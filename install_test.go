package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestZipDetection(t *testing.T) {
	data, err := ioutil.ReadFile("test_data/test.zip")
	if err != nil {
		t.Fatalf("error reading test.zip: %v", err)
	}

	if !isZipFile(data) {
		t.Fatal("isZipFile(<test.zip>); want true, got false")
	}

	data, err = ioutil.ReadFile("test_data/OFL.txt")
	if err != nil {
		t.Fatalf("error reading OFL.txt: %v", err)
	}

	if isZipFile(data) {
		t.Fatal("isZipFile(<OFL.txt>); want false, got true")
	}
}

func TestFontInstallation(t *testing.T) {
	filename := "test_data/open-sans-v18-latin-regular.ttf"
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Fatalf("error reading %v: %v", filename, err)
	}

	fontData, err := NewFontData(filename, data)
	if err != nil {
		t.Fatalf(`NewFontData(%s, ..) should not return an error`, filename)
	}

	// It seems silly to test what is basically a simple file copy but...here we are.
	if err = install(fontData, FontsDir); err != nil {
		t.Fatalf("error installing font to %v: %v", FontsDir, err)
	}

	dest := filepath.Join(FontsDir, filepath.Base(filename))
	if _, err = os.Stat(dest); err != nil {
		t.Fatalf("installed font does not exist in %v", FontsDir)
	} else {
		if err := os.Remove(dest); err != nil {
			t.Fatalf("error removing installed font: %v", dest)
		}
	}
}

func TestFontInstallationFromZip(t *testing.T) {
	zip := "test_data/test.zip"
	data, err := ioutil.ReadFile(zip)
	if err != nil {
		t.Fatalf("error reading %v: %v", zip, err)
	}

	if err = installFromZIP(data, FontsDir); err != nil {
		t.Fatalf("error installing font to %v: %v", FontsDir, err)
	}

	files := map[string]bool{
		"open-sans-v18-latin-regular.ttf": true,
		"Apache License.txt":              false,
	}

	for filename, shouldExist := range files {
		dest := filepath.Join(FontsDir, filename)
		if shouldExist {
			if _, err = os.Stat(dest); err != nil {
				t.Fatalf("file %v does not exist in %v", filename, FontsDir)
			} else {
				if err := os.Remove(dest); err != nil {
					t.Fatalf("error removing installed font: %v", dest)
				}
			}
		} else {
			if _, err = os.Stat(dest); err == nil {
				t.Fatalf("file %v exists in %v but shouldn't", filename, FontsDir)
			}
		}
	}
}
