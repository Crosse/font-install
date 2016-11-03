package main

import (
	"path"

	xdg "github.com/casimir/xdg-go"
)

// FontsDir denotes the path to the user's fonts directory on Linux.
var FontsDir = path.Join(xdg.DataHome(), "fonts")
