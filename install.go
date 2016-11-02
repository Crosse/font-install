// +build darwin linux

package main

import (
	"archive/zip"
	"bytes"
	"errors"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/ConradIrwin/font/sfnt"
	log "github.com/Crosse/gosimplelogger"
)

func (f Font) install(compressedFile *zip.File) (err error) {
	rc, err := compressedFile.Open()
	if err != nil {
		return
	}
	defer rc.Close()

	buf, err := ioutil.ReadAll(rc)
	if err != nil {
		return
	}

	rs := bytes.NewReader(buf)

	if _, ok := FontExtensions[path.Ext(compressedFile.Name)]; !ok {
		// Only install files that are actual fonts.
		log.Debugf("Non-font file not installed: \"%v\"", compressedFile.Name)
		return
	}

	font, err := sfnt.Parse(rs)
	if err != nil {
		return
	}

	if font.HasTable(sfnt.TagName) == false {
		return errors.New("Font has no name table")
	}

	nameTable := font.NameTable()
	entries := make(map[sfnt.NameID]string)
	for _, nameEntry := range nameTable.List() {
		entries[nameEntry.NameID] = nameEntry.String()
	}
	name := entries[sfnt.NameFull]
	family := entries[sfnt.NamePreferredFamily]
	if family == "" {
		if v, ok := entries[sfnt.NameFontFamily]; ok {
			family = v
		} else {
			log.Errorf("Font %v has no font family!", name)
		}
	}

	if name == "" {
		log.Errorf("Font %v has no name!", compressedFile.Name)
		name = compressedFile.Name
	}
	log.Infof("Installing %v", name)

	fontPath := path.Join(FontsDir, strings.ToLower(strings.Replace(family, " ", "-", -1)))
	log.Debugf("Creating font directory %v", fontPath)
	err = os.MkdirAll(fontPath, 0700)
	if err != nil {
		return
	}

	fileName := path.Join(fontPath, compressedFile.Name)
	if err = os.MkdirAll(path.Dir(fileName), 0700); err != nil {
		return
	}

	if err = ioutil.WriteFile(fileName, buf, 0644); err != nil {
		return
	}

	return nil
}
