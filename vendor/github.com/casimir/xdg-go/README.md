xdg-go [![GoDoc](https://godoc.org/github.com/casimir/xdg-go?status.svg)](https://godoc.org/github.com/casimir/xdg-go) [![codebeat badge](https://codebeat.co/badges/845ce4ee-6285-45dc-a790-e56c00d0f35c)](https://codebeat.co/projects/github-com-casimir-xdg-go)
=======================================================================================================================================================================================================================================

## Quickstart

If you just want OS-sensible paths.
```
configDirs := xdg.ConfigDirs()
dataPath := xdg.DataHome()
cachePath := xdg.CacheHome()
```

Alternatively you can create a context that would determine full paths for your application files.
```
app := xdg.App{Name: "someApp"}
configFile := app.ConfigPath("someApp.toml")
dataFile := app.DataPath("data.json")
```

## Supported path types

This is a KISS implementation of the XDG Base Directory Specification. As of now it handles the following path types:
- Data (`XDG_DATA_*`) for application-wide or user-wide data.
- Config (`XDG_CONFIG_*`) for application-wide or user-wide config.
- Cache (`XDG_CACHE_*`)for application-wide or user-wide cached data.

## Multi-OS

The specification is Linux centric but this implementation targets more: Linux, OSX and Windows. Default values has been chosen regarding both the specification and the OS conventions. Note than you can override these values with the corresponding environment variables.

There are a lot of OSes missing but supporting them implies a good knowledge of these conventions and philosophies, contributors maybe?


