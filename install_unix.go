// +build linux solaris openbsd

package main

import (
	"io/ioutil"
	"os"
	"path"
	"strings"

	log "github.com/Crosse/gosimplelogger"
)

func platformDependentInstall(fontData *FontData) (err error) {
	// On Linux, fontconfig can understand subdirectories. So, to keep the
	// font directory clean, install all font files for a particular font
	// family into a subdirectory named after the family (with hyphens instead
	// of spaces).
	fullPath := path.Join(FontsDir,
		strings.ToLower(strings.Replace(fontData.Family, " ", "-", -1)),
		path.Base(fontData.FileName))
	log.Debugf("Installing \"%v\" to %v", fontData.Name, fullPath)

	if err = os.MkdirAll(path.Dir(fullPath), 0700); err != nil {
		return err
	}

	err = ioutil.WriteFile(fullPath, fontData.Data, 0644) //nolint:gosec

	return nil
}
