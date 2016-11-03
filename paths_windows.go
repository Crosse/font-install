package main

import (
	"os"
	"path"
)

// FontsDir denotes the path to the user's fonts directory on Linux.
// Windows doesn't have the concept of a permanent, per-user collection
// of fonts, meaning that all fonts are stored in the system-level fonts
// directory, which is %WINDIR%\Fonts by default.
var FontsDir = path.Join(os.Getenv("WINDIR"), "Fonts")
