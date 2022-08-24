//go:build linux || solaris || openbsd || freebsd
// +build linux solaris openbsd freebsd

package main

import (
	"fmt"
	"os"
	"path"
	"strings"

	log "github.com/Crosse/gosimplelogger"
)

func platformDependentInstall(fontData *FontData) error {
	// On Linux, fontconfig can understand subdirectories. So, to keep the
	// font directory clean, install all font files for a particular font
	// family into a subdirectory named after the family (with hyphens instead
	// of spaces).
	fullPath := path.Join(FontsDir,
		strings.ToLower(strings.ReplaceAll(fontData.Family, " ", "-")),
		path.Base(fontData.FileName))
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
