package main

import (
	"os"
	"path/filepath"
)

// FontsDir is the current user's Windows font directory.
var FontsDir = filepath.Join(os.Getenv("LOCALAPPDATA"), "Microsoft", "Windows", "Fonts")

// GlobalFontsDir is the system-wide Windows font directory.
var GlobalFontsDir = filepath.Join(os.Getenv("WINDIR"), "Fonts")
