package main

import (
	"archive/zip"
	"bytes"
	"fmt"
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
	var b []byte
	var fontData *FontData

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

	if isZipFile(b) {
		err = installFromZIP(b)
	} else {
		fontData, err = NewFontData(path.Base(u.Path), b)
		if err != nil {
			return
		}
		err = install(fontData)
	}

	return
}

func isZipFile(data []byte) bool {
	contentType := http.DetectContentType(data)
	log.Debugf("Detected content type: %v", contentType)

	return contentType == "application/zip"
}

func getRemoteFile(url string) (data []byte, err error) {
	log.Infof("Downloading font file from %v", url)

	var client = http.Client{}

	resp, err := client.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
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

func installFromZIP(data []byte) (err error) {
	bytesReader := bytes.NewReader(data)

	zipReader, err := zip.NewReader(bytesReader, int64(bytesReader.Len()))
	if err != nil {
		return
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

		fontData, err := NewFontData(zf.Name, data)
		if err != nil {
			log.Errorf(`Skipping non-font file "%s"`, zf.Name)
			continue
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

	return
}

func install(fontData *FontData) (err error) {
	log.Infof("==> %s", fontData.Name)

	err = platformDependentInstall(fontData)
	if err == nil {
		installedFonts += 1
	}

	return
}
