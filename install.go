package main

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"runtime"
	"strings"

	log "github.com/Crosse/gosimplelogger"
)

var installedFonts = 0

// InstallFont installs the font specified by fontPath.
// fontPath can either be a URL or a filesystem path.
// For URLs, only the "file", "http", and "https" schemes are currently valid.
func InstallFont(fontPath string) (err error) {
	var (
		b        []byte
		fontData *FontData
	)

	u, err := url.Parse(fontPath)
	if err != nil {
		return
	}

	switch u.Scheme {
	case "file", "":
		if b, err = getLocalFile(fontPath); err != nil {
			return err
		}
	case "http", "https":
		if b, err = getRemoteFile(fontPath); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unhandled URL scheme: %v", u.Scheme)
	}

	filename := path.Base(u.Path)

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
		fontData, err = NewFontData(filename, b)
		if err != nil {
			return err
		}

		return install(fontData)
	}

}

func getContentType(data []byte) string {
	contentType := http.DetectContentType(data)
	log.Debugf("Detected content type: %v", contentType)

	return contentType
}

func getRemoteFile(url string) (data []byte, err error) {
	log.Infof("Downloading font file from %v", url)

	var client = http.Client{}

	resp, err := client.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Debugf("HTTP request resulted in status %v", resp.StatusCode)
		return
	}

	data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	return
}

func getLocalFile(filename string) (data []byte, err error) {
	if data, err = ioutil.ReadFile(filename); err != nil {
		return nil, err
	}

	return
}

func installFromGZIP(filename string, data []byte) (err error) {
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

func installFromTarball(r io.Reader) (err error) {
	log.Debug("reading tarball")

	tarReader := tar.NewReader(r)

	fonts := make(map[string]*FontData)

	log.Debug("Scanning tarball for fonts")

	for {
		hdr, err := tarReader.Next()
		if err == io.EOF {
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

func installFromZIP(data []byte) (err error) {
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
			return err
		}
		defer rc.Close()

		data, err := ioutil.ReadAll(rc)
		if err != nil {
			return err
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

func installFonts(fonts map[string]*FontData) (err error) {
	for _, font := range fonts {
		if strings.Contains(strings.ToLower(font.Name), "windows compatible") {
			if runtime.GOOS != "windows" {
				// hack to not install the "Windows Compatible" version of every nerd font.
				log.Infof(`Ignoring "%s" on non-Windows platform`, font.Name)
				continue
			}
		}

		if err = install(font); err != nil {
			return err
		}
	}

	return nil
}

func install(fontData *FontData) (err error) {
	log.Infof("==> %s", fontData.Name)

	err = platformDependentInstall(fontData)
	if err == nil {
		installedFonts++
	}

	return err
}
