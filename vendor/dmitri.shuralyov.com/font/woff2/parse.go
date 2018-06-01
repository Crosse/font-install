package woff2

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/dsnet/compress/brotli"
)

// File represents a parsed WOFF2 file.
type File struct {
	Header         Header
	TableDirectory TableDirectory
	// CollectionDirectory is present only if the font is a collection,
	// as reported by Header.IsCollection.
	CollectionDirectory *CollectionDirectory

	// FontData is the concatenation of data for each table in the font.
	// During storage, it's compressed using Brotli.
	FontData []byte

	ExtendedMetadata *ExtendedMetadata

	// PrivateData is an optional block of private data for the font designer,
	// foundry, or vendor to use.
	PrivateData []byte
}

// Parse parses the WOFF2 data from r.
func Parse(r io.Reader) (File, error) {
	hdr, err := parseHeader(r)
	if err != nil {
		return File{}, err
	}
	td, err := parseTableDirectory(r, hdr)
	if err != nil {
		return File{}, err
	}
	cd, err := parseCollectionDirectory(r, hdr)
	if err != nil {
		return File{}, err
	}
	fd, err := parseCompressedFontData(r, hdr, td)
	if err != nil {
		return File{}, err
	}
	em, err := parseExtendedMetadata(r, hdr)
	if err != nil {
		return File{}, err
	}
	pd, err := parsePrivateData(r, hdr)
	if err != nil {
		return File{}, err
	}

	// Check for padding with a maximum of three null bytes.
	// TODO: This check needs to be moved to Extended Metadata and Private Data blocks,
	//       and made more precise (i.e., the beginning of those blocks must be 4-byte aligned).
	n, err := io.Copy(discardZeroes{}, r)
	if err != nil {
		return File{}, fmt.Errorf("Parse: %v", err)
	}
	if n > 3 {
		return File{}, fmt.Errorf("Parse: %d bytes left remaining, want no more than 3", n)
	}

	return File{
		Header:              hdr,
		TableDirectory:      td,
		CollectionDirectory: cd,
		FontData:            fd,
		ExtendedMetadata:    em,
		PrivateData:         pd,
	}, nil
}

// discardZeroes is an io.Writer that returns an error if any non-zero bytes are written to it.
type discardZeroes struct{}

func (discardZeroes) Write(p []byte) (int, error) {
	for _, b := range p {
		if b != 0 {
			return 0, fmt.Errorf("encountered non-zero byte %d", b)
		}
	}
	return len(p), nil
}

// Header is the file header with basic font type and version,
// along with offsets to metadata and private data blocks.
type Header struct {
	Signature           uint32 // The identifying signature; must be 0x774F4632 ('wOF2').
	Flavor              uint32 // The "sfnt version" of the input font.
	Length              uint32 // Total size of the WOFF file.
	NumTables           uint16 // Number of entries in directory of font tables.
	Reserved            uint16 // Reserved; set to 0.
	TotalSfntSize       uint32 // Total size needed for the uncompressed font data, including the sfnt header, directory, and font tables (including padding).
	TotalCompressedSize uint32 // Total length of the compressed data block.
	MajorVersion        uint16 // Major version of the WOFF file.
	MinorVersion        uint16 // Minor version of the WOFF file.
	MetaOffset          uint32 // Offset to metadata block, from beginning of WOFF file.
	MetaLength          uint32 // Length of compressed metadata block.
	MetaOrigLength      uint32 // Uncompressed size of metadata block.
	PrivOffset          uint32 // Offset to private data block, from beginning of WOFF file.
	PrivLength          uint32 // Length of private data block.
}

func parseHeader(r io.Reader) (Header, error) {
	var hdr Header
	err := binary.Read(r, order, &hdr)
	if err != nil {
		return Header{}, err
	}
	if hdr.Signature != signature {
		return Header{}, fmt.Errorf("parseHeader: invalid signature: got %#08x, want %#08x", hdr.Signature, signature)
	}
	return hdr, nil
}

// IsCollection reports whether this is a font collection, i.e.,
// if the value of Flavor field is set to the TrueType Collection flavor 'ttcf'.
func (hdr Header) IsCollection() bool {
	return hdr.Flavor == ttcfFlavor
}

// TableDirectory is the directory of font tables, containing size and other info.
type TableDirectory []TableDirectoryEntry

func parseTableDirectory(r io.Reader, hdr Header) (TableDirectory, error) {
	var td TableDirectory
	for i := 0; i < int(hdr.NumTables); i++ {
		var e TableDirectoryEntry

		err := readU8(r, &e.Flags)
		if err != nil {
			return nil, err
		}
		if e.Flags&0x3f == 0x3f {
			e.Tag = new(uint32)
			err := readU32(r, e.Tag)
			if err != nil {
				return nil, err
			}
		}
		err = readBase128(r, &e.OrigLength)
		if err != nil {
			return nil, err
		}

		switch tag, transformVersion := e.tag(), e.transformVersion(); tag {
		case glyfTable, locaTable:
			// 0 means transform for glyf/loca tables.
			if transformVersion == 0 {
				e.TransformLength = new(uint32)
				err := readBase128(r, e.TransformLength)
				if err != nil {
					return nil, err
				}

				// The transform length of the transformed loca table MUST always be zero.
				if tag == locaTable && *e.TransformLength != 0 {
					return nil, fmt.Errorf("parseTableDirectory: 'loca' table has non-zero transform length %d", *e.TransformLength)
				}
			}
		default:
			// Non-0 means transform for other tables.
			if transformVersion != 0 {
				e.TransformLength = new(uint32)
				err := readBase128(r, e.TransformLength)
				if err != nil {
					return nil, err
				}
			}
		}

		td = append(td, e)
	}
	return td, nil
}

// Table is a high-level representation of a table.
type Table struct {
	Tag    uint32
	Offset int
	Length int
}

// Tables returns the derived high-level information
// about the tables in the table directory.
func (td TableDirectory) Tables() []Table {
	var ts []Table
	var offset int
	for _, t := range td {
		length := int(t.length())
		ts = append(ts, Table{
			Tag:    t.tag(),
			Offset: offset,
			Length: length,
		})
		offset += length
	}
	return ts
}

// uncompressedSize computes the total uncompressed size
// of the tables in the table directory.
func (td TableDirectory) uncompressedSize() int64 {
	var n int64
	for _, t := range td {
		n += int64(t.length())
	}
	return n
}

// TableDirectoryEntry is a table directory entry.
type TableDirectoryEntry struct {
	Flags           uint8   // Table type and flags.
	Tag             *uint32 // 4-byte tag (optional).
	OrigLength      uint32  // Length of original table.
	TransformLength *uint32 // Transformed length (optional).
}

func (e TableDirectoryEntry) tag() uint32 {
	switch e.Tag {
	case nil:
		return knownTableTags[e.Flags&0x3f] // Bits [0..5].
	default:
		return *e.Tag
	}
}

func (e TableDirectoryEntry) transformVersion() uint8 {
	return e.Flags >> 6 // Bits [6..7].
}

func (e TableDirectoryEntry) length() uint32 {
	switch e.TransformLength {
	case nil:
		return e.OrigLength
	default:
		return *e.TransformLength
	}
}

// CollectionDirectory is an optional table containing the font fragment descriptions
// of font collection entries.
type CollectionDirectory struct {
	Header  CollectionHeader
	Entries []CollectionFontEntry
}

// CollectionHeader is a part of CollectionDirectory.
type CollectionHeader struct {
	Version  uint32
	NumFonts uint16
}

// CollectionFontEntry represents a CollectionFontEntry record.
type CollectionFontEntry struct {
	NumTables    uint16   // The number of tables in this font.
	Flavor       uint32   // The "sfnt version" of the font.
	TableIndices []uint16 // The indicies identifying an entry in the Table Directory for each table in this font.
}

func parseCollectionDirectory(r io.Reader, hdr Header) (*CollectionDirectory, error) {
	// CollectionDirectory is present only if the input font is a collection.
	if !hdr.IsCollection() {
		return nil, nil
	}

	var cd CollectionDirectory
	err := readU32(r, &cd.Header.Version)
	if err != nil {
		return nil, err
	}
	err = read255UShort(r, &cd.Header.NumFonts)
	if err != nil {
		return nil, err
	}
	for i := 0; i < int(cd.Header.NumFonts); i++ {
		var e CollectionFontEntry

		err := read255UShort(r, &e.NumTables)
		if err != nil {
			return nil, err
		}
		err = readU32(r, &e.Flavor)
		if err != nil {
			return nil, err
		}
		for j := 0; j < int(e.NumTables); j++ {
			var tableIndex uint16
			err := read255UShort(r, &tableIndex)
			if err != nil {
				return nil, err
			}
			if tableIndex >= hdr.NumTables {
				return nil, fmt.Errorf("parseCollectionDirectory: tableIndex >= hdr.NumTables")
			}
			e.TableIndices = append(e.TableIndices, tableIndex)
		}

		cd.Entries = append(cd.Entries, e)
	}
	return &cd, nil
}

func parseCompressedFontData(r io.Reader, hdr Header, td TableDirectory) ([]byte, error) {
	// Compressed font data.
	br, err := brotli.NewReader(io.LimitReader(r, int64(hdr.TotalCompressedSize)), nil)
	//br, err := brotli.NewReader(&exactReader{R: r, N: int64(hdr.TotalCompressedSize)}, nil)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	n, err := io.Copy(&buf, br)
	if err != nil {
		return nil, fmt.Errorf("parseCompressedFontData: io.Copy: %v", err)
	}
	err = br.Close()
	if err != nil {
		return nil, fmt.Errorf("parseCompressedFontData: br.Close: %v", err)
	}
	if uncompressedSize := td.uncompressedSize(); n != uncompressedSize {
		return nil, fmt.Errorf("parseCompressedFontData: unexpected size of uncompressed data: got %d, want %d", n, uncompressedSize)
	}
	return buf.Bytes(), nil
}

// ExtendedMetadata is an optional block of extended metadata,
// represented in XML format and compressed for storage in the WOFF2 file.
type ExtendedMetadata struct{}

func parseExtendedMetadata(r io.Reader, hdr Header) (*ExtendedMetadata, error) {
	if hdr.MetaLength == 0 {
		return nil, nil
	}
	return nil, fmt.Errorf("parseExtendedMetadata: not implemented")
}

func parsePrivateData(r io.Reader, hdr Header) ([]byte, error) {
	if hdr.PrivLength == 0 {
		return nil, nil
	}
	return nil, fmt.Errorf("parsePrivateData: not implemented")
}

// readU8 reads a UInt8 value.
func readU8(r io.Reader, v *uint8) error {
	return binary.Read(r, order, v)
}

// readU16 reads a UInt16 value.
func readU16(r io.Reader, v *uint16) error {
	return binary.Read(r, order, v)
}

// readU32 reads a UInt32 value.
func readU32(r io.Reader, v *uint32) error {
	return binary.Read(r, order, v)
}

// readBase128 reads a UIntBase128 value.
func readBase128(r io.Reader, v *uint32) error {
	var accum uint32
	for i := 0; i < 5; i++ {
		var data uint8
		err := binary.Read(r, order, &data)
		if err != nil {
			return err
		}

		// Leading zeros are invalid.
		if i == 0 && data == 0x80 {
			return fmt.Errorf("leading zero is invalid")
		}

		// If any of top 7 bits are set then accum << 7 would overflow.
		if accum&0xfe000000 != 0 {
			return fmt.Errorf("top seven bits are set, about to overflow")
		}

		accum = (accum << 7) | uint32(data)&0x7f

		// Spin until most significant bit of data byte is false.
		if (data & 0x80) == 0 {
			*v = accum
			return nil
		}
	}
	return fmt.Errorf("UIntBase128 sequence exceeds 5 bytes")
}

// read255UShort reads a 255UInt16 value.
func read255UShort(r io.Reader, v *uint16) error {
	const (
		oneMoreByteCode1 = 255
		oneMoreByteCode2 = 254
		wordCode         = 253
		lowestUCode      = 253
	)
	var code uint8
	err := binary.Read(r, order, &code)
	if err != nil {
		return err
	}
	switch code {
	case wordCode:
		var value uint16
		err := binary.Read(r, order, &value)
		if err != nil {
			return err
		}
		*v = value
		return nil
	case oneMoreByteCode1:
		var value uint8
		err := binary.Read(r, order, &value)
		if err != nil {
			return err
		}
		*v = uint16(value) + lowestUCode
		return nil
	case oneMoreByteCode2:
		var value uint8
		err := binary.Read(r, order, &value)
		if err != nil {
			return err
		}
		*v = uint16(value) + lowestUCode*2
		return nil
	default:
		*v = uint16(code)
		return nil
	}
}

// WOFF2 uses big endian encoding.
var order binary.ByteOrder = binary.BigEndian
