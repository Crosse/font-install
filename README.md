# font-install

[![Build Status](https://travis-ci.org/Crosse/font-install.svg?branch=master)](https://travis-ci.org/Crosse/font-install)

`font-install` is a cross-platform utility to install fonts on a system.
It can install fonts on Linux, OSX, or Windows systems.
Given a ZIP file, it will even install all font files within the archive.
If you feed `font-install` an HTTP/HTTPS URL, it will first download the file, then install the font (or extract and install them if the file is a ZIP file).
`font-install` currently only handles OpenType and TrueType font files.

`font-install` is not intended to handle webfonts; it installs fonts so that other applications can use (such as your display manager, office suite, etc.).
It currently installs fonts into the system's user-specific fonts library location.

## Background

I have a list of fonts that I always want installed on my computer, no matter which operating system that computer runs.
On Linux, this evolved into a simple bash script, which also worked well for OSX (after fudging the install path).
Both of these operating systems simply look for font files in a specific location (on Linux, `${HOME}/.local/share/fonts`; for OSX it is `${HOME}/Library/Fonts`).
However, I also wanted to be able to install these same fonts just as easily on Windows...which is quite not as easy. (It's still pretty easy, though.)

Enter `font-install`: a tool that will download, extract, and install fonts the exact same way no matter which OS I'm on.

(Finally, as with all my personal projects, a big reason for building this tool was to also work on my development skills.)

## Requirements

This tool was originally built with Go 1.7.3, and has been verified to build with 1.10.2.

## Installation

```
$ go install github.com/Crosse/font-install
```

## Usage

General usage details:
```
font-install [-debug] [-fromFile] [font ...]
    -debug                  enable debug logging
    -fromFile <string>      text file containing fonts to install
```

## Examples
* Download and install the [Source Sans Pro][source-sans-pro] font from [Font Squirrel][fontsquirrel]:
    ```
    $ font-install http://www.fontsquirrel.com/fonts/download/source-sans-pro

    Installing Source Sans Pro ExtraLight (SourceSansPro-ExtraLight.otf)
    Installing Source Sans Pro ExtraLight Italic (SourceSansPro-ExtraLightIt.otf)
    Installing Source Sans Pro Light (SourceSansPro-Light.otf)
    Installing Source Sans Pro Light Italic (SourceSansPro-LightIt.otf)
    Installing Source Sans Pro (SourceSansPro-Regular.otf)
    Installing Source Sans Pro Italic (SourceSansPro-It.otf)
    Installing Source Sans Pro Semibold (SourceSansPro-Semibold.otf)
    Installing Source Sans Pro Semibold Italic (SourceSansPro-SemiboldIt.otf)
    Installing Source Sans Pro Bold (SourceSansPro-Bold.otf)
    Installing Source Sans Pro Bold Italic (SourceSansPro-BoldIt.otf)
    Installing Source Sans Pro Black (SourceSansPro-Black.otf)
    Installing Source Sans Pro Black Italic (SourceSansPro-BlackIt.otf)
    ```
  This downloads the source-sans-pro ZIP archive from Font Squirrel, extracts all of the fonts, and installs them into the user's fonts directory.

* Download and install [Chopin Script][chopin-script] from [dafont.com] and enable debug output:
    ```
    $ font-install -debug 'http://dl.dafont.com/dl/?f=chopin_script'

    Installing font from http://dl.dafont.com/dl/?f=chopin_script
    Downloading font file from http://dl.dafont.com/dl/?f=chopin_script
    Detected content type: application/zip
    Installing ChopinScript
    Installing "ChopinScript" to /Users/seth/Library/Fonts/ChopinScript.otf
    ```

* Install a font file you have stored locally on your file system:
    ```
    $ file *.ttf
    DroidSans-Bold.ttf: TrueType font data
    DroidSans.ttf:      TrueType font data
    DroidSansMono.ttf:  TrueType font data

    $ font-install *.ttf
    Installing Droid Sans Bold
    Installing Droid Sans
    Installing Droid Sans Mono
    ```
* Feed `font-install` a list of fonts to install via a [text file][example.txt]:
    ```
    $ font-install -fromFile example.txt

    Installing Open Sans Light
    Installing Open Sans Light Italic
    Installing Open Sans
    Installing Open Sans Italic
    Installing Open Sans Semibold
    Installing Open Sans Semibold Italic
    Installing Open Sans Bold
    Installing Open Sans Bold Italic
    Installing Open Sans Extrabold
    Installing Open Sans Extrabold Italic
    Installing Source Code Pro ExtraLight
    Installing Source Code Pro ExtraLight Italic
    Installing Source Code Pro Light
    [...]
    Installing FontAwesome
    ```


[chopin-script]: http://www.dafont.com/chopin-script.font
[dafont.com]: http://www.dafont.com
[example.txt]: example.txt
[fontsquirrel]: https://www.fontsquirrel.com
[source-sans-pro]: https://www.fontsquirrel.com/fonts/source-sans-pro
