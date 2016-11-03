package main

import (
	"io"
	"os"
	"path"
	"strings"
)

func platformDependentInstall(fontData *FontData) (err error) {
	// On Linux, fontconfig can understand subdirectories. So, to keep the
	// font directory clean, install all font files for a particular font
	// family into a subdirectory named after the family (with hyphens instead
	// of spaces).
	fullPath := path.Join(FontsDir, strings.ToLower(strings.Replace(fontData.Family, " ", "-", -1), fontData.FileName))
	log.Debugf("Installing \"%v\" to %v", fontData.Name, fullPath)

	if err = os.MkdirAll(path.Dir(fullPath), 0700); err != nil {
		return
	}

	fd, err := os.Create(fullPath)
	if err != nil {
		return
	}
	defer fd.Close()

	_, err = io.Copy(fd, fontData.Data)
	return
}
