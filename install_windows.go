package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"syscall"
	"unsafe"

	log "github.com/Crosse/gosimplelogger"
	"golang.org/x/sys/windows/registry"
)

func platformDependentInstall(fontData *FontData) error {
	// Windows needs the file, a registry entry, and a font-resource notification.
	fileName := path.Base(fontData.FileName)
	fullPath := filepath.Join(FontsDir, fileName)
	fontPath, err := syscall.UTF16PtrFromString(fullPath)
	if err != nil {
		return fmt.Errorf("invalid font path: %w", err)
	}
	log.Debugf("Installing \"%v\" to %v", fontData.Name, fullPath)

	if !globalInstall {
		if err := os.MkdirAll(FontsDir, 0o700); err != nil {
			return fmt.Errorf("cannot create font directory: %w", err)
		}
	}
	err = os.WriteFile(fullPath, fontData.Data, 0o644)
	if err != nil {
		return fmt.Errorf("cannot write file: %w", err)
	}

	// Per-user font entries need an absolute path; system entries use the filename.
	hive := registry.CURRENT_USER
	registryValue := fullPath
	if globalInstall {
		hive = registry.LOCAL_MACHINE
		registryValue = fileName
	}
	k, _, err := registry.CreateKey(
		hive,
		`SOFTWARE\Microsoft\Windows NT\CurrentVersion\Fonts`,
		registry.SET_VALUE)
	if err != nil {
		// If this fails, remove the font file as well.
		log.Error(err)

		if nexterr := os.Remove(fullPath); nexterr != nil {
			return fmt.Errorf("error removing font: %w", nexterr)
		}

		return fmt.Errorf("error opening registry: %w", err)
	}
	defer k.Close()

	valueName := fmt.Sprintf("%s (TrueType)", fontData.Name)
	if err = k.SetStringValue(valueName, registryValue); err != nil {
		// If this fails, remove the font file as well.
		log.Error(err)

		if nexterr := os.Remove(fullPath); nexterr != nil {
			return fmt.Errorf("error removing font: %w", nexterr)
		}

		return fmt.Errorf("error writing to registry: %w", err)
	}

	added, _, callErr := syscall.NewLazyDLL("gdi32.dll").NewProc("AddFontResourceW").Call(uintptr(unsafe.Pointer(fontPath)))
	if added == 0 {
		if err := k.DeleteValue(valueName); err != nil {
			log.Error(err)
		}
		if err := os.Remove(fullPath); err != nil {
			log.Error(err)
		}
		return fmt.Errorf("cannot load installed font %s: %v", fullPath, callErr)
	}

	// Tell running applications to refresh their font lists without a logoff.
	const hwndBroadcast, wmFontChange = 0xffff, 0x001d
	syscall.NewLazyDLL("user32.dll").NewProc("SendNotifyMessageW").Call(hwndBroadcast, wmFontChange, 0, 0)
	return nil
}
