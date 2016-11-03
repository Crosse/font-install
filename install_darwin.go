package main

import (
	"io"
	"os"
	"path"

	log "github.com/Crosse/gosimplelogger"
)

func platformDependentInstall(fontData *FontData) (err error) {
	// On darwin/OSX, the user's fonts directory is ~/Library/Fonts,
	// and fonts should be installed directly into that path;
	// i.e., not in subfolders.
	fullPath := path.Join(FontsDir, path.Base(fontData.FileName))
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
