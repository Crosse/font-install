package main

import (
	"fmt"
	"os"
	"path"

	log "github.com/Crosse/gosimplelogger"
)

func platformDependentInstall(fontData *FontData) error {
	// On darwin/OSX, the user's fonts directory is ~/Library/Fonts,
	// and fonts should be installed directly into that path;
	// i.e., not in subfolders.
	fullPath := path.Join(FontsDir, path.Base(fontData.FileName))
	log.Debugf("Installing \"%v\" to %v", fontData.Name, fullPath)

	err := os.MkdirAll(path.Dir(fullPath), 0o700)
	if err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	err = os.WriteFile(fullPath, fontData.Data, 0o644)
	if err != nil {
		return fmt.Errorf("cannot write file: %w", err)
	}

	return nil
}
