package sfnt

import (
	"encoding/binary"
	"fmt"
	"io"
)

type woffHeader struct {
	Signature      Tag
	Flavor         Tag
	Length         uint32
	NumTables      uint16
	Reserved       uint16
	TotalSfntSize  uint32
	Version        fixed
	MetaOffset     uint32
	MetaLength     uint32
	MetaOrigLength uint32
	PrivOffset     uint32
	PrivLength     uint32
}

type woffEntry struct {
	Tag          Tag
	Offset       uint32
	CompLength   uint32
	OrigLength   uint32
	OrigChecksum uint32
}

func readWOFFHeader(r io.Reader, header *woffHeader) error {
	return binary.Read(r, binary.BigEndian, header)
}

func readWOFFHeaderFast(r io.Reader, header *woffHeader) error {
	var buf [44]byte
	if _, err := io.ReadFull(r, buf[:]); err != nil {
		return err
	}

	header.Signature = NewTag(buf[0:4])
	header.Flavor = NewTag(buf[4:8])
	header.Length = binary.BigEndian.Uint32(buf[8:12])
	header.NumTables = binary.BigEndian.Uint16(buf[12:14])
	header.Reserved = binary.BigEndian.Uint16(buf[14:16])
	header.TotalSfntSize = binary.BigEndian.Uint32(buf[16:20])
	header.Version.Major = int16(binary.BigEndian.Uint16(buf[20:22]))
	header.Version.Minor = binary.BigEndian.Uint16(buf[22:24])
	header.MetaOffset = binary.BigEndian.Uint32(buf[24:28])
	header.MetaLength = binary.BigEndian.Uint32(buf[28:32])
	header.MetaOrigLength = binary.BigEndian.Uint32(buf[32:36])
	header.PrivOffset = binary.BigEndian.Uint32(buf[36:40])
	header.PrivLength = binary.BigEndian.Uint32(buf[40:44])
	return nil
}

func readWOFFEntry(r io.Reader, entry *woffEntry) error {
	return binary.Read(r, binary.BigEndian, entry)
}

func readWOFFEntryFast(r io.Reader, entry *woffEntry) error {
	var buf [20]byte
	if _, err := io.ReadFull(r, buf[:]); err != nil {
		return err
	}
	entry.Tag = NewTag(buf[0:4])
	entry.Offset = binary.BigEndian.Uint32(buf[4:8])
	entry.CompLength = binary.BigEndian.Uint32(buf[8:12])
	entry.OrigLength = binary.BigEndian.Uint32(buf[12:16])
	entry.OrigChecksum = binary.BigEndian.Uint32(buf[16:20])
	return nil
}

func parseWOFF(file File) (*Font, error) {
	var header woffHeader
	if err := readWOFFHeaderFast(file, &header); err != nil {
		return nil, err
	}

	font := &Font{
		file:       file,
		scalerType: header.Flavor,
		tables:     make(map[Tag]*tableSection, header.NumTables),
	}

	for i := 0; i < int(header.NumTables); i++ {
		var entry woffEntry
		if err := readWOFFEntryFast(file, &entry); err != nil {
			return nil, err
		}

		// TODO Check the checksum.

		if _, found := font.tables[entry.Tag]; found {
			return nil, fmt.Errorf("found multiple %q tables", entry.Tag)
		}

		font.tables[entry.Tag] = &tableSection{
			tag: entry.Tag,

			offset:  entry.Offset,
			length:  entry.CompLength,
			zLength: entry.OrigLength,
		}
	}

	if _, ok := font.tables[TagHead]; !ok {
		return nil, ErrMissingHead
	}

	return font, nil
}
