package main

import (
	"archive/zip"
	"bytes"
	"io/ioutil"
	"net/http"

	log "github.com/Crosse/gosimplelogger"
)

var FontExtensions = map[string]bool{
	".otf": true,
	".ttf": true,
}

type Font struct {
	Name string
	URL  string
}
type Fonts []Font

func (f Font) Download() (zipReader *zip.Reader, err error) {
	zipReader = nil
	err = nil

	log.Debugf("Downloading font ZIP file from %v", f.URL)

	var client = http.Client{}
	resp, err := client.Get(f.URL)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Debugf("HTTP request resulted in status %v", resp.StatusCode)
		return
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	r := bytes.NewReader(b)
	zipReader, err = zip.NewReader(r, int64(len(b)))
	if err != nil {
		return
	}

	return
}

func (f Font) Install() error {
	if zipReader, err := f.Download(); err == nil {
		for _, zf := range zipReader.File {
			f.install(zf)
		}
	} else {
		return err
	}
	return nil
}
