package main

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"

	log "github.com/Crosse/gosimplelogger"
)

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
		return errors.New(fmt.Sprintf("Unhandled URL scheme: %v", u.Scheme))
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
	log.Debugf("Downloading font file from %v", url)

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

		if fontData, err := NewFontData(zf.Name, data); err == nil {
			err = install(fontData)
		}
	}
	return
}

func install(fontData *FontData) (err error) {
	log.Infof("Installing %v (%v)", fontData.Name, fontData.FileName)
	return platformDependentInstall(fontData)
}
