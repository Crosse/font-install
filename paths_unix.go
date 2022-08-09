// +build linux solaris openbsd freebsd

package main

import (
	"os"
	"path"
)

// FontsDir denotes the path to the user's fonts directory on Unix-like systems.
var FontsDir = path.Join(os.Getenv("HOME"), "/.local/share/fonts")
