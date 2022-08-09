# font-install

`font-install` is a cross-platform utility to install fonts on a system.

It can install fonts on Linux, macOS, OpenBSD, FreeBSD, or Windows systems.
Given a ZIP file, it will even install all font files within the archive.
If you feed `font-install` an HTTP/HTTPS URL, it will first download the
file, then install the font (or extract and install them if the file is a
ZIP file).  `font-install` currently only handles OpenType and TrueType font
files.

`font-install` is not intended to handle webfonts; it installs fonts so that
other applications can use (such as your display manager, office suite,
etc.).  It currently installs fonts into the system's user-specific fonts
library location.

## Background

I have a list of fonts that I always want installed on my computer, no
matter which operating system that computer runs.  On Linux, this evolved
into a simple bash script, which also worked well for OSX (after fudging the
install path).  Both of these operating systems simply look for font files
in a specific location (on Linux, `${HOME}/.local/share/fonts`; for OSX it
is `${HOME}/Library/Fonts`).  However, I also wanted to be able to install
these same fonts just as easily on Windows...which is quite not as
easy. (It's still pretty easy, though.)

Enter `font-install`: a tool that will download, extract, and install fonts
the exact same way no matter which OS I'm on.

## Requirements

This code requires at least Go 1.18 to build.

## Installation

You can find the latest release on the [releases page][releases] for many
platforms. For other platforms, you may install from source like so:

```
$ go install github.com/Crosse/font-install@latest
```

## Usage

General usage details:
```
Usage of font-install:
  -debug
        Enable debug logging
  -dry-run
        Don't actually download or install anything
  -fromFile string
        text file containing fonts to install
```

## Examples
* Download and install the [Source Sans Pro][source-sans-pro] font from
  [Font Squirrel][fontsquirrel]:
   ```
   $ font-install http://www.fontsquirrel.com/fonts/download/source-sans-pro
   Downloading font file from http://www.fontsquirrel.com/fonts/download/source-sans-pro
   Skipping non-font file "SIL Open Font License.txt"
   ==> Source Sans Pro Black
   ==> Source Sans Pro Black Italic
   ==> Source Sans Pro Light
   ==> Source Sans Pro Light Italic
   ==> Source Sans Pro Italic
   ==> Source Sans Pro Semibold
   ==> Source Sans Pro Semibold Italic
   ==> Source Sans Pro ExtraLight
   ==> Source Sans Pro ExtraLight Italic
   ==> Source Sans Pro
   ==> Source Sans Pro Bold
   ==> Source Sans Pro Bold Italic
   Installed 12 fonts
   ```

  This downloads the source-sans-pro ZIP archive from Font Squirrel,
  extracts all of the fonts, and installs them into the user's fonts
  directory.

* Download and install [Chopin Script][chopin-script] from [dafont.com] and enable debug output:
  ```
  $ font-install -debug 'http://dl.dafont.com/dl/?f=chopin_script'
  Installing font from http://dl.dafont.com/dl/?f=chopin_script
  Downloading font file from http://dl.dafont.com/dl/?f=chopin_script
  Detected content type: application/zip
  Scanning ZIP file for fonts
  ==> ChopinScript
  Installing "ChopinScript" to /Users/seth/Library/Fonts/ChopinScript.otf
  Installed 1 fonts
  ```

* Install a font file you have stored locally on your file system:
  ```
  $ ls *.ttf
  OpenSans-Bold.ttf    OpenSans-Italic.ttf  OpenSans-Regular.ttf

  $ font-install *.ttf
  ==> Open Sans Bold
  ==> Open Sans Italic
  ==> Open Sans
  Installed 3 fonts
  ```
* Feed `font-install` a list of fonts to install via a [text file][example.txt]:
  ```
  $ font-install -fromFile example.txt
  Downloading font file from http://www.fontsquirrel.com/fonts/download/Inconsolata
  Skipping non-font file "SIL Open Font License.txt"
  Downloading font file from http://www.fontsquirrel.com/fonts/download/dejavu-sans
  Skipping non-font file "DejaVu Fonts License.txt"
  ==> DejaVu Sans Oblique
  ==> DejaVu Sans Condensed Oblique
  ==> DejaVu Sans ExtraLight
  ==> DejaVu Sans
  ==> DejaVu Sans Condensed
  ==> DejaVu Sans Condensed Bold
  ==> DejaVu Sans Condensed Bold Oblique
  ==> DejaVu Sans Bold
  ==> DejaVu Sans Bold Oblique
  Downloading font file from http://www.fontsquirrel.com/fonts/download/dejavu-sans-mono
  Skipping non-font file "DejaVu Fonts License.txt"
  ==> DejaVu Sans Mono Bold Oblique
  ==> DejaVu Sans Mono
  ==> DejaVu Sans Mono Oblique
  ==> DejaVu Sans Mono Bold
  [...]
  Downloading font file from http://www.fontsquirrel.com/fonts/download/ubuntu-mono
  Skipping non-font file "UBUNTU FONT LICENCE.txt"
  ==> Ubuntu Mono
  ==> Ubuntu Mono Italic
  ==> Ubuntu Mono Bold
  ==> Ubuntu Mono Bold Italic
  Downloading font file from http://fontawesome.io/assets/font-awesome-4.7.0.zip
  not a font: font-awesome-4.7.0.zip
  Installed 72 fonts
  ```

  The above output also shows how `font-install` handles errors. (In this
  case, the URL is incorrect.)


[releases]: https://github.com/Crosse/font-install/releases/latest
[chopin-script]: http://www.dafont.com/chopin-script.font
[dafont.com]: http://www.dafont.com
[example.txt]: example.txt
[fontsquirrel]: https://www.fontsquirrel.com
[source-sans-pro]: https://www.fontsquirrel.com/fonts/source-sans-pro
