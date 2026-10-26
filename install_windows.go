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

var (
	procAddFontResourceW    = syscall.NewLazyDLL("gdi32.dll").NewProc("AddFontResourceW")
	procRemoveFontResourceW = syscall.NewLazyDLL("gdi32.dll").NewProc("RemoveFontResourceW")
	procSendNotifyMessageW  = syscall.NewLazyDLL("user32.dll").NewProc("SendNotifyMessageW")
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
	// A loaded font keeps its file open, so unload every reference to an
	// existing copy before overwriting it. The same number of references is
	// restored afterward so applications still using the font keep working.
	unloaded := 0
	for {
		removed, _, _ := procRemoveFontResourceW.Call(uintptr(unsafe.Pointer(fontPath)))
		if removed == 0 {
			break
		}
		unloaded++
	}
	reload := func(n int) {
		for i := 0; i < n; i++ {
			procAddFontResourceW.Call(uintptr(unsafe.Pointer(fontPath)))
		}
	}
	err = os.WriteFile(fullPath, fontData.Data, 0o644)
	if err != nil {
		reload(unloaded)
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

	added, _, callErr := procAddFontResourceW.Call(uintptr(unsafe.Pointer(fontPath)))
	if added == 0 {
		if err := k.DeleteValue(valueName); err != nil {
			log.Error(err)
		}
		if err := os.Remove(fullPath); err != nil {
			log.Error(err)
		}
		return fmt.Errorf("cannot load installed font %s: %v", fullPath, callErr)
	}
	reload(unloaded - 1)

	// Tell running applications to refresh their font lists without a logoff.
	const hwndBroadcast, wmFontChange = 0xffff, 0x001d
	procSendNotifyMessageW.Call(hwndBroadcast, wmFontChange, 0, 0)
	return nil
}
