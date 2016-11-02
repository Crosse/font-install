package main

import (
	"path"

	xdg "github.com/casimir/xdg-go"
)

var FontsDir = path.Join(xdg.DataHome(), "Fonts")
