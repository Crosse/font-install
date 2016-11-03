// +build darwin linux

package main

import (
	"io"
	"os"
	"path"
	"strings"
)

func platformDependentInstall(fontData *FontData) (err error) {
	fontPath := path.Join(FontsDir, strings.ToLower(strings.Replace(fontData.Family, " ", "-", -1)))
	err = os.MkdirAll(fontPath, 0700)
	if err != nil {
		return
	}

	fullPath := path.Join(fontPath, fontData.FileName)
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
