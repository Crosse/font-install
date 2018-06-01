font
====

[![Build Status](https://travis-ci.org/ConradIrwin/font.svg?branch=master)](https://travis-ci.org/ConradIrwin/font) [![GoDoc](https://godoc.org/github.com/ConradIrwin/font?status.svg)](https://godoc.org/github.com/ConradIrwin/font)

A collection of Go packages for parsing and encoding OpenType fonts.

The main contribution of this repository is the [SFNT](https://godoc.org/github.com/ConradIrwin/font/sfnt) library which provides support for parsing OpenType, TrueType, WOFF, and WOFF2 fonts.

Also included is a utility called `font` that can do various useful things with fonts:

```
go get -u github.com/ConradIrwin/font/cmd/font
```

Info gets information about the font from the `name` table:

```
font info ~/Downloads/Fanwood.ttf
```

Scrub empties the name table (which can give you a few kb savings, even if you gzip or woff2-encode your font).

```
font scrub ~/Downloads/Fanwood.ttf
```

Stats tells you how much space each table is using:

```
font stats ~/Downloads/Fanwood.ttf
```

TODO
----

Still missing is support for parsing EOT files (which should be easy to add). Also support for generating WOFF files (which is annoyingly fiddly due to the checksum calculation) and WOFF2 files (needs a Brotli encoder), and a whole load of code around dealing with the hundreds of other SFNT table formats.

Font file formats
-----------------

On the web there are four main types of font file, TrueType, OpenType, WOFF, WOFF2, and EOT. They all represent the same SFNT information, but are encoded slightly differently. You may also come across SVG fonts, which are a totally different beast.

Inside one of these files, there are two main types of glyphs, TrueType and
OpenType (also known as PostScript Type 2, or CFF). There are also a series of supporting
tables which contain meta-data about the font (its Name, Copyright Information, Kerning tweaks, Ligatures, etc.etc.)

To confuse things further, OpenType fonts use exactly the same format as TrueType fonts, and a WOFF file can contain an OpenType glyphs or a TrueType glyphs. There's no good solution to the ambiguity in terminolgy, just be aware of it.

License
-------

Copyright (c) Conrad Irwin 2015, MIT license. See LICENSE.MIT for details
