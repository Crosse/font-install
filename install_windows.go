package main

import (
	"archive/zip"
	"bytes"
	"errors"
	"io/ioutil"
	"path"

	"github.com/ConradIrwin/font/sfnt"
	log "github.com/Crosse/gosimplelogger"
	"golang.org/x/sys/windows/registry"
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

	baseName := path.Base(compressedFile.Name)
	if _, ok := FontExtensions[path.Ext(baseName)]; !ok {
		// Only install files that are actual fonts.
		log.Debugf("Non-font file not installed: \"%v\"", baseName)
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
		log.Errorf("Font %v has no name!", baseName)
		name = baseName
	}
	log.Infof("Installing %v", name)

	// To install a font on Windows:
	//  - Copy the file to the fonts directory
	//  - Create a registry entry for the font
	fileName := path.Join(FontsDir, baseName)
	if err = ioutil.WriteFile(fileName, buf, 0644); err != nil {
		return
	}

	k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows NT\CurrentVersion\Fonts`, registry.WRITE)
	if err != nil {
		return
	}
	defer k.Close()
	if err = k.SetStringValue(name, baseName); err != nil {
		//TODO: clean up the file we just created?
		return
	}

	return nil
}
