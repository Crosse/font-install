package sfnt

import (
	"bytes"
	"encoding/binary"
	"io"
)

type tableOS2Fields struct {
	Version             uint16
	XAvgCharWidth       uint16
	USWeightClass       uint16
	USWidthClass        uint16
	FSType              uint16
	YSubscriptXSize     int16
	YSubscriptYSize     int16
	YSubscriptXOffset   int16
	YSubscriptYOffset   int16
	YSuperscriptXSize   int16
	YSuperscriptYSize   int16
	YSuperscriptXOffset int16
	YSuperscriptYOffset int16
	YStrikeoutSize      int16
	YStrikeoutPosition  int16
	SFamilyClass        int16
	Panose              [10]byte
	UlCharRange         [4]uint32
	AchVendID           Tag
	FsSelection         uint16
	FsFirstCharIndex    uint16
	FsLastCharIndex     uint16
	STypoAscender       int16
	STypoDescender      int16
	STypoLineGap        int16
	UsWinAscent         uint16
	UsWinDescent        uint16
	UlCodePageRange1    uint32
	UlCodePageRange2    uint32
	SxHeigh             int16
	SCapHeight          int16
	UsDefaultChar       uint16
	UsBreakChar         uint16
	UsMaxContext        uint16
	UsLowerPointSize    uint16
	UsUpperPointSize    uint16
}

type TableOS2 struct {
	baseTable
	tableOS2Fields
	bytes []byte
}

func parseTableOS2(tag Tag, buf []byte) (Table, error) {
	r := bytes.NewBuffer(buf)

	var table tableOS2Fields
	if err := binary.Read(r, binary.BigEndian, &table); err != nil {
		// Different versions of the table are different lengths, as such
		// we may not already read every field.
		if err != io.ErrUnexpectedEOF {
			return nil, err
		}

		// TODO Check the len(buf) is expected for this version
	}

	return &TableOS2{
		baseTable:      baseTable(tag),
		tableOS2Fields: table,
		bytes:          buf,
	}, nil
}

func (t *TableOS2) Bytes() []byte {
	return t.bytes
}
