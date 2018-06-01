package sfnt

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
)

var (
	// TagHead represents the 'head' table, which contains the font header
	TagHead = MustNamedTag("head")
	// TagMaxp represents the 'maxp' table, which contains the maximum profile
	TagMaxp = MustNamedTag("maxp")
	// TagHmtx represents the 'hmtx' table, which contains the horizontal metrics
	TagHmtx = MustNamedTag("hmtx")
	// TagHhea represents the 'hhea' table, which contains the horizonal header
	TagHhea = MustNamedTag("hhea")
	// TagOS2 represents the 'OS/2' table, which contains windows-specific metadata
	TagOS2 = MustNamedTag("OS/2")
	// TagName represents the 'name' table, which contains font name information
	TagName = MustNamedTag("name")
	// TagGpos represents the 'GPOS' table, which contains Glyph Positioning features
	TagGpos = MustNamedTag("GPOS")
	// TagGsub represents the 'GSUB' table, which contains Glyph Substitution features
	TagGsub = MustNamedTag("GSUB")

	// TypeTrueType is the first four bytes of an OpenType file containing a TrueType font
	TypeTrueType = Tag{0x00010000}
	// TypeAppleTrueType is the first four bytes of an OpenType file containing a TrueType font
	// (specifically one designed for Apple products, it's recommended to use TypeTrueType instead)
	TypeAppleTrueType = MustNamedTag("true")
	// TypePostScript1 is the first four bytes of an OpenType file containing a PostScript Type 1 font
	TypePostScript1 = MustNamedTag("typ1")
	// TypeOpenType is the first four bytes of an OpenType file containing a PostScript Type 2 font
	// as specified by OpenType
	TypeOpenType = MustNamedTag("OTTO")

	// SignatureWOFF is the magic number at the start of a WOFF file.
	SignatureWOFF = MustNamedTag("wOFF")

	// SignatureWOFF2 is the magic number at the start of a WOFF2 file.
	SignatureWOFF2 = MustNamedTag("wOF2")
)

// Tag represents an open-type table name.
// These are technically uint32's, but are usually
// displayed in ASCII as they are all acronyms.
// see https://developer.apple.com/fonts/TrueType-Reference-Manual/RM06/Chap6.html#Overview
type Tag struct {
	Number uint32
}

// NamedTag gives you the Tag corresponding to the acronym.
func NamedTag(str string) (Tag, error) {
	bytes := []byte(str)

	if len(bytes) != 4 {
		return Tag{}, fmt.Errorf("invalid tag: must be exactly 4 bytes")
	}

	return Tag{uint32(bytes[0])<<24 |
		uint32(bytes[1])<<16 |
		uint32(bytes[2])<<8 |
		uint32(bytes[3])}, nil
}

// MustNamedTag gives you the Tag corresponding to the acronym.
// This function will panic if the string passed in is not 4 bytes long.
func MustNamedTag(str string) Tag {
	t, err := NamedTag(str)
	if err != nil {
		panic(err)
	}
	return t
}

func NewTag(bytes []byte) Tag {
	return Tag{Number: binary.BigEndian.Uint32(bytes)}
}

func ReadTag(r io.Reader) (Tag, error) {
	bytes := make([]byte, 4)
	_, err := io.ReadFull(r, bytes)
	return NewTag(bytes), err
}

// String returns the ASCII representation of the tag.
func (tag Tag) String() string {
	return string(tag.bytes())
}

func (tag Tag) bytes() []byte {
	return []byte{
		byte(tag.Number >> 24 & 0xFF),
		byte(tag.Number >> 16 & 0xFF),
		byte(tag.Number >> 8 & 0xFF),
		byte(tag.Number & 0xFF),
	}
}

func (tag Tag) hex() string {
	return "0x" + hex.EncodeToString(tag.bytes())
}
