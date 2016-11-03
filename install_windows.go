package main

import (
	"io"
	"os"
	"path"

	log "github.com/Crosse/gosimplelogger"
	"golang.org/x/sys/windows/registry"
)

func platformDependentInstall(fontData *FontData) (err error) {
	// To install a font on Windows:
	//  - Copy the file to the fonts directory
	//  - Create a registry entry for the font
	fullPath := path.Join(FontsDir, fontData.FileName)

	fd, err := os.Create(fullPath)
	if err != nil {
		return
	}
	defer fd.Close()

	_, err = io.Copy(fd, fontData.Data)
	if err != nil {
		return
	}

	k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows NT\CurrentVersion\Fonts`, registry.WRITE)
	if err != nil {
		log.Error(err)
		if nexterr := os.Remove(fullPath); nexterr != nil {
			return nexterr
		}
		return err
	}
	defer k.Close()
	if err = k.SetStringValue(fontData.Name, fontData.FileName); err != nil {
		log.Error(err)
		if nexterr := os.Remove(fullPath); nexterr != nil {
			return nexterr
		}
		return
	}

	return nil
}
