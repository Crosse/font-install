package sfnt

import (
	"bytes"
	"encoding/binary"
)

// TableHead contains critical information about the rest of the font.
// https://developer.apple.com/fonts/TrueType-Reference-Manual/RM06/Chap6head.html
type TableHead struct {
	baseTable
	tableHeadFields
}

type tableHeadFields struct {
	VersionNumber      fixed
	FontRevision       fixed
	CheckSumAdjustment uint32
	MagicNumber        uint32
	Flags              uint16
	UnitsPerEm         uint16
	Created            longdatetime
	Updated            longdatetime
	XMin               int16
	YMin               int16
	XMax               int16
	YMax               int16
	MacStyle           uint16
	LowestRecPPEM      uint16
	FontDirection      int16
	IndexToLocFormat   int16
	GlyphDataFormat    int16
}

func parseTableHead(tag Tag, buf []byte) (Table, error) {
	r := bytes.NewBuffer(buf)

	var fields tableHeadFields
	if err := binary.Read(r, binary.BigEndian, &fields); err != nil {
		return nil, err
	}

	return &TableHead{
		baseTable:       baseTable(tag),
		tableHeadFields: fields,
	}, nil
}

// Bytes returns the byte representation of this header.
func (table *TableHead) Bytes() []byte {
	var buffer bytes.Buffer
	if err := binary.Write(&buffer, binary.BigEndian, table); err != nil {
		panic(err) // should never happen
	}
	return buffer.Bytes()
}

// ExpectedChecksum is the checksum that the file should have had.
func (table *TableHead) ExpectedChecksum() uint32 {
	return 0xB1B0AFBA - table.CheckSumAdjustment
}

// SetExpectedChecksum updates the table so it can be output with the correct checksum.
func (table *TableHead) SetExpectedChecksum(checksum uint32) {
	table.CheckSumAdjustment = 0xB1B0AFBA - checksum
}

// ClearExpectedChecksum updates the table so that the checksum can be calculated.
func (table *TableHead) ClearExpectedChecksum() {
	table.CheckSumAdjustment = 0
}
