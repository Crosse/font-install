package main

import (
	"path"

	xdg "github.com/casimir/xdg-go"
)

// FontsDir denotes the path to the user's fonts directory on OSX.
var FontsDir = path.Join(xdg.DataHome(), "Fonts")
