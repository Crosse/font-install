package main

import "os"

// The path to the Fonts directory on Windows.
// Windows doesn't have the concept of a permanent, per-user collection
// of fonts, meaning that all fonts are stored in the system-level fonts
// directory, which is %WINDIR%\Fonts by default.
const FontsDir = paths.Join(os.Getenv("WINDIR"), "Fonts")
