package main

import (
	"fmt"
	"os"
	"path"

	log "github.com/Crosse/gosimplelogger"
	"golang.org/x/sys/windows/registry"
)

func platformDependentInstall(fontData *FontData) error {
	// To install a font on Windows:
	//  - Copy the file to the fonts directory
	//  - Create a registry entry for the font
	fullPath := path.Join(FontsDir, fontData.FileName)
	log.Debugf("Installing \"%v\" to %v", fontData.Name, fullPath)

	err := os.WriteFile(fullPath, fontData.Data, 0o644)
	if err != nil {
		return fmt.Errorf("cannot write file: %w", err)
	}

	// Second, write metadata about the font to the registry.
	k, err := registry.OpenKey(
		registry.LOCAL_MACHINE,
		`SOFTWARE\Microsoft\Windows NT\CurrentVersion\Fonts`,
		registry.WRITE)
	if err != nil {
		// If this fails, remove the font file as well.
		log.Error(err)

		if nexterr := os.Remove(fullPath); nexterr != nil {
			return fmt.Errorf("error removing font: %w", nexterr)
		}

		return fmt.Errorf("error opening registry: %w", err)
	}
	defer k.Close()

	// Apparently it's "ok" to mark an OpenType font as "TrueType",
	// and since this tool only supports True- and OpenType fonts,
	// this should be Okay(tm).
	// Besides, Windows does it, so why can't I?
	valueName := fmt.Sprintf("%v (TrueType)", fontData.FileName)
	if err = k.SetStringValue(fontData.Name, valueName); err != nil {
		// If this fails, remove the font file as well.
		log.Error(err)

		if nexterr := os.Remove(fullPath); nexterr != nil {
			return fmt.Errorf("error removing font: %w", nexterr)
		}

		return fmt.Errorf("error writing to registry: %w", err)
	}

	return nil
}
