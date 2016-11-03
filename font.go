package main

import (
	"bytes"
	"errors"
	"fmt"
	"path"

	"github.com/ConradIrwin/font/sfnt"
	log "github.com/Crosse/gosimplelogger"
)

type FontData struct {
	Name     string
	Family   string
	FileName string
	Metadata map[sfnt.NameID]string
	Data     *bytes.Reader
}

var FontExtensions = map[string]bool{
	".otf": true,
	".ttf": true,
}

func NewFontData(fileName string, data []byte) (fontData *FontData, err error) {
	if _, ok := FontExtensions[path.Ext(fileName)]; !ok {
		return nil, errors.New(fmt.Sprintf("Not a font: %v", fileName))
	}

	fontData = &FontData{
		FileName: fileName,
		Metadata: make(map[sfnt.NameID]string),
		Data:     bytes.NewReader(data),
	}

	font, err := sfnt.Parse(fontData.Data)
	if err != nil {
		return nil, err
	}

	if font.HasTable(sfnt.TagName) == false {
		return nil, errors.New(fmt.Sprintf("Font %v has no name table", fileName))
	}

	nameTable := font.NameTable()
	for _, nameEntry := range nameTable.List() {
		fontData.Metadata[nameEntry.NameID] = nameEntry.String()
	}
	fontData.Name = fontData.Metadata[sfnt.NameFull]
	fontData.Family = fontData.Metadata[sfnt.NamePreferredFamily]
	if fontData.Family == "" {
		if v, ok := fontData.Metadata[sfnt.NameFontFamily]; ok {
			fontData.Family = v
		} else {
			log.Errorf("Font %v has no font family!", fontData.Name)
		}
	}

	if fontData.Name == "" {
		log.Errorf("Font %v has no name! Using file name instead.", fileName)
		fontData.Name = fileName
	}

	return
}
