package main

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"syscall"
	"unsafe"

	log "github.com/Crosse/gosimplelogger"
	"golang.org/x/sys/windows"
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
	pending := false
	if globalInstall && fileInUse(err) {
		// Fonts in %WINDIR%\Fonts can stay mapped by other processes even after
		// unloading. Like Windows Installer, replace the file at the next boot.
		err = replaceAtReboot(fullPath, fontData.Data)
		pending = err == nil
	}
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
		registry.QUERY_VALUE|registry.SET_VALUE)
	if err != nil {
		// If this fails, remove the font file as well.
		log.Error(err)

		if pending {
			reload(unloaded)
		} else if nexterr := os.Remove(fullPath); nexterr != nil {
			return fmt.Errorf("error removing font: %w", nexterr)
		}

		return fmt.Errorf("error opening registry: %w", err)
	}
	defer k.Close()

	fontType := "TrueType"
	if strings.EqualFold(filepath.Ext(fileName), ".otf") {
		fontType = "OpenType"
	}
	valueName := fmt.Sprintf("%s (%s)", fontData.Name, fontType)
	prevValue, _, prevErr := k.GetStringValue(valueName)
	if err = k.SetStringValue(valueName, registryValue); err != nil {
		// If this fails, remove the font file as well.
		log.Error(err)

		if pending {
			reload(unloaded)
		} else if nexterr := os.Remove(fullPath); nexterr != nil {
			return fmt.Errorf("error removing font: %w", nexterr)
		}

		return fmt.Errorf("error writing to registry: %w", err)
	}

	if pending {
		// The old file stays in place until the reboot, so keep it loaded.
		reload(unloaded)
		log.Infof("%s is in use and will be updated after a reboot", fullPath)
		return nil
	}

	added, _, callErr := procAddFontResourceW.Call(uintptr(unsafe.Pointer(fontPath)))
	if added == 0 {
		// Restore an earlier registration of this font if its file is still
		// there; otherwise remove the entry that was just written.
		var regErr error
		if prevErr == nil && otherFontExists(prevValue, fullPath) {
			regErr = k.SetStringValue(valueName, prevValue)
		} else {
			regErr = k.DeleteValue(valueName)
		}
		if regErr != nil {
			log.Error(regErr)
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

// otherFontExists reports whether a font registry value names an existing
// file other than fullPath.
func otherFontExists(value, fullPath string) bool {
	prevPath := value
	if !filepath.IsAbs(prevPath) {
		prevPath = filepath.Join(FontsDir, prevPath)
	}
	if strings.EqualFold(filepath.Clean(prevPath), filepath.Clean(fullPath)) {
		return false
	}
	_, err := os.Stat(prevPath)
	return err == nil
}

func fileInUse(err error) bool {
	return errors.Is(err, windows.ERROR_SHARING_VIOLATION) ||
		errors.Is(err, windows.ERROR_USER_MAPPED_FILE)
}

// replaceAtReboot writes data next to fullPath and schedules it to replace
// fullPath at the next boot. This requires administrator privileges.
func replaceAtReboot(fullPath string, data []byte) error {
	tmp, err := os.CreateTemp(filepath.Dir(fullPath), filepath.Base(fullPath)+".*.tmp")
	if err != nil {
		return err
	}
	_, err = tmp.Write(data)
	if closeErr := tmp.Close(); err == nil {
		err = closeErr
	}

	var from, to *uint16
	if err == nil {
		from, err = windows.UTF16PtrFromString(tmp.Name())
	}
	if err == nil {
		to, err = windows.UTF16PtrFromString(fullPath)
	}
	if err == nil {
		err = windows.MoveFileEx(from, to,
			windows.MOVEFILE_REPLACE_EXISTING|windows.MOVEFILE_DELAY_UNTIL_REBOOT)
	}
	if err != nil {
		os.Remove(tmp.Name())
		return fmt.Errorf("cannot schedule font replacement: %w", err)
	}

	return nil
}
