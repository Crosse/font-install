package main

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	log "github.com/Crosse/gosimplelogger"
)

var installedFonts = 0

// InstallFont installs the font specified by fontPath.
// fontPath can either be a URL or a filesystem path.
// For URLs, only the "file", "http", and "https" schemes are currently valid.
func InstallFont(fontPath string) error {
	return installFontPath(fontPath, true)
}

func installFontPath(fontPath string, expandGlob bool) error {
	// Only arguments that look like URLs are parsed as URLs. Parsing plain
	// paths misreads Windows drive letters and colons as schemes, and percent
	// signs as escapes.
	if strings.Contains(fontPath, "://") {
		return installFromURL(fontPath)
	}

	// Shells on Windows do not expand wildcards for command arguments.
	// An existing file is used as-is, since font names often contain
	// brackets (e.g. "Inter[opsz,wght].ttf").
	_, statErr := os.Stat(fontPath)
	if expandGlob && statErr != nil && strings.ContainsAny(fontPath, "*?[") {
		return installGlob(fontPath)
	}

	b, err := getLocalFile(fontPath)
	if err != nil {
		return err
	}

	return installData(filepath.Base(fontPath), b)
}

// installFromURL installs the font at a file, http, or https URL.
func installFromURL(rawURL string) error {
	u, err := url.Parse(rawURL)
	if err != nil {
		return fmt.Errorf("error parsing URL: %w", err)
	}

	switch u.Scheme {
	case "file":
		localPath, err := fileURLToPath(u)
		if err != nil {
			return err
		}

		return installFontPath(localPath, false)
	case "http", "https":
		b, err := getRemoteFile(rawURL)
		if err != nil {
			return err
		}

		return installData(path.Base(u.Path), b)
	default:
		return fmt.Errorf("unhandled URL scheme: %v", u.Scheme)
	}
}

// fileURLToPath converts a file URL to a local path. On Windows, a URL with
// a host becomes a UNC path.
func fileURLToPath(u *url.URL) (string, error) {
	hasHost := u.Host != "" && !strings.EqualFold(u.Host, "localhost")

	if runtime.GOOS != "windows" {
		if hasHost {
			return "", fmt.Errorf("unsupported file URL host: %q", u.Host)
		}

		return u.Path, nil
	}

	// file:///C:/fonts/x.ttf has the path "/C:/fonts/x.ttf".
	p := u.Path
	if len(p) >= 3 && p[0] == '/' && p[2] == ':' {
		p = p[1:]
	}

	p = filepath.FromSlash(p)
	if hasHost {
		p = `\\` + u.Host + p
	}

	return p, nil
}

// installData installs the font, or the fonts in the archive, in b.
// filename is used to identify the file type.
func installData(filename string, b []byte) error {
	ct := getContentType(b)
	log.Debugf("content type: %s", ct)

	switch ct {
	case "application/zip":
		return installFromZIP(b)
	case "application/x-gzip":
		return installFromGZIP(filename, b)
	case "application/octet-stream":
		if strings.ToLower(path.Ext(filename)) == ".tar" {
			return installFromTarball(bytes.NewReader(b))
		}

		fallthrough
	default:
		fontData, err := NewFontData(filename, b)
		if err != nil {
			return err
		}

		return install(fontData)
	}
}

// installGlob installs every local file matching pattern. Matches are
// installed as literal paths, so brackets in their names aren't expanded again.
func installGlob(pattern string) error {
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return fmt.Errorf("invalid file pattern: %w", err)
	}

	if len(matches) == 0 {
		return fmt.Errorf("no files match %q", pattern)
	}

	for _, match := range matches {
		if err := installFontPath(match, false); err != nil {
			return err
		}
	}

	return nil
}

func getContentType(data []byte) string {
	contentType := http.DetectContentType(data)
	log.Debugf("Detected content type: %v", contentType)

	return contentType
}

func getRemoteFile(url string) ([]byte, error) {
	log.Infof("Downloading font file from %v", url)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("cannot make http request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error getting remote file: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Debugf("HTTP request resulted in status %v", resp.StatusCode)
		return nil, fmt.Errorf("server returned non-successful status code %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("erorr reading remote file: %w", err)
	}

	return data, nil
}

func getLocalFile(filename string) ([]byte, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("cannot read local file: %w", err)
	}

	return data, nil
}

func installFromGZIP(filename string, data []byte) error {
	log.Debug("reading gzipped file")

	bytesReader := bytes.NewReader(data)

	gzipReader, err := gzip.NewReader(bytesReader)
	if err != nil {
		return fmt.Errorf("cannot read gzip file: %w", err)
	}
	defer gzipReader.Close()

	uncompressedFilename := strings.TrimSuffix(filename, ".gz")
	ext := strings.ToLower(path.Ext(uncompressedFilename))

	if ext == ".tar" || ext == ".tgz" {
		return installFromTarball(gzipReader)
	}

	// Gzipped files only contain a single compressed file, so we'll just assume that it's one compressed font.
	b, err := io.ReadAll(gzipReader)
	if err != nil {
		return fmt.Errorf("cannot read compressed file: %w", err)
	}

	fontData, err := NewFontData(path.Base(uncompressedFilename), b)
	if err != nil {
		return err
	}

	return install(fontData)
}

func installFromTarball(r io.Reader) error {
	log.Debug("reading tarball")

	tarReader := tar.NewReader(r)

	fonts := make(map[string]*FontData)

	log.Debug("Scanning tarball for fonts")

	for {
		hdr, err := tarReader.Next()
		if errors.Is(err, io.EOF) {
			break
		}

		if err != nil {
			return fmt.Errorf("cannot read tarball: %w", err)
		}

		data, err := io.ReadAll(tarReader)
		if err != nil {
			return fmt.Errorf("unable to read file %s from tarball: %w", hdr.Name, err)
		}

		appendFont(fonts, hdr.Name, data)
	}

	return installFonts(fonts)
}

func installFromZIP(data []byte) error {
	log.Debug("reading zipfile")

	bytesReader := bytes.NewReader(data)

	zipReader, err := zip.NewReader(bytesReader, int64(bytesReader.Len()))
	if err != nil {
		return fmt.Errorf("cannot read zip file: %w", err)
	}

	fonts := make(map[string]*FontData)

	log.Debug("Scanning ZIP file for fonts")

	for _, zf := range zipReader.File {
		rc, err := zf.Open()
		if err != nil {
			return fmt.Errorf("cannot open compressed file %s: %w", zf.Name, err)
		}
		defer rc.Close()

		data, err := io.ReadAll(rc)
		if err != nil {
			return fmt.Errorf("cannot read compressed file %s: %w", zf.Name, err)
		}

		appendFont(fonts, zf.Name, data)
	}

	return installFonts(fonts)
}

func appendFont(fonts map[string]*FontData, fileName string, data []byte) {
	fontData, err := NewFontData(fileName, data)
	if err != nil {
		log.Errorf(`Skipping non-font file "%s"`, fileName)
		return
	}

	if _, ok := fonts[fontData.Name]; !ok {
		fonts[fontData.Name] = fontData
	} else {
		// Prefer OTF over TTF; otherwise prefer the first font we found.
		first := strings.ToLower(path.Ext(fonts[fontData.Name].FileName))
		second := strings.ToLower(path.Ext(fontData.FileName))

		if first != second && second == ".otf" {
			log.Infof(`Preferring "%s" over "%s"`, fontData.FileName, fonts[fontData.Name].FileName)
			fonts[fontData.Name] = fontData
		}
	}
}

func installFonts(fonts map[string]*FontData) error {
	for _, font := range fonts {
		if strings.Contains(strings.ToLower(font.Name), "windows compatible") {
			if runtime.GOOS != "windows" {
				// hack to not install the "Windows Compatible" version of every nerd font.
				log.Infof(`Ignoring "%s" on non-Windows platform`, font.Name)
				continue
			}
		}

		if err := install(font); err != nil {
			return err
		}
	}

	return nil
}

func install(fontData *FontData) error {
	log.Infof("==> %s", fontData.Name)

	err := platformDependentInstall(fontData)
	if err == nil {
		installedFonts++
	}

	return err
}
