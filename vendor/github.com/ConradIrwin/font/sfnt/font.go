package sfnt

import (
	"errors"
	"fmt"
	"sort"
)

type fixed struct {
	Major int16
	Minor uint16
}

type longdatetime struct {
	SecondsSince1904 uint64
}

func (u *unparsedTable) Bytes() []byte {
	return u.bytes
}

// ErrMissingHead is returned by ParseOTF when the font has no head section.
var ErrMissingHead = errors.New("missing head table in font")

// ErrInvalidChecksum is returned by ParseOTF if the font's checksum is wrong
var ErrInvalidChecksum = errors.New("invalid checksum")

// ErrUnsupportedFormat is returned from Parse if parsing failed
var ErrUnsupportedFormat = errors.New("unsupported font format")

// ErrMissingTable is returned from *Table if the table does not exist in the font.
var ErrMissingTable = errors.New("missing table")

// Font represents a SFNT font, which is the underlying representation found
// in .otf and .ttf files (and .woff, .woff2, .eot files)
// SFNT is a container format, which contains a number of tables identified by
// Tags. Depending on the type of glyphs embedded in the file which tables will
// exist. In particular, there's a big different between TrueType glyphs (usually .ttf)
// and CFF/PostScript Type 2 glyphs (usually .otf)
type Font struct {
	file File

	scalerType Tag
	tables     map[Tag]*tableSection
}

// tableSection represents a table within the font file.
type tableSection struct {
	tag   Tag
	table Table

	offset  uint32 // Offset into the file this table starts.
	length  uint32 // Length of this table within the file.
	zLength uint32 // Uncompressed length of this table.
}

// Tags is the list of tags that are defined in this font, sorted by numeric value.
func (font *Font) Tags() []Tag {
	tags := make([]Tag, 0, len(font.tables))

	for t := range font.tables {
		tags = append(tags, t)
	}

	sort.Slice(tags, func(i, j int) bool {
		return tags[i].Number < tags[j].Number
	})

	return tags
}

// HasTable returns true if this font has an entry for the given table.
func (font *Font) HasTable(tag Tag) bool {
	_, ok := font.tables[tag]
	return ok
}

// AddTable adds a table to the font. If a table with the
// given tag is already present, it will be overwritten.
func (font *Font) AddTable(tag Tag, table Table) {
	font.tables[tag] = &tableSection{
		tag:   tag,
		table: table,
	}
}

// RemoveTable removes a table from the font. If the table
// doesn't exist, this method will do nothing.
func (font *Font) RemoveTable(tag Tag) {
	delete(font.tables, tag)
}

// Type represents the kind of glyphs in this font.
// It is one of TypeTrueType, TypeTrueTypeApple, TypePostScript1, TypeOpenType
func (font *Font) Type() Tag {
	return font.scalerType
}

// String provides a debugging representation of a font.
func (font *Font) String() string {
	str := "Parsed font with scalerType=" + font.scalerType.hex()

	if font.scalerType != TypeTrueType {
		str += " (" + font.scalerType.String() + ")"
	}

	for _, t := range font.Tags() {
		str += "\n" + t.String()
	}

	return str
}

// HeadTable returns the table corresponding to the 'head' tag.
func (font *Font) HeadTable() (*TableHead, error) {
	t, err := font.Table(TagHead)
	if err != nil {
		return nil, err
	}
	return t.(*TableHead), nil
}

// NameTable returns the table corresponding to the 'name' tag.
func (font *Font) NameTable() (*TableName, error) {
	t, err := font.Table(TagName)
	if err != nil {
		return nil, err
	}
	return t.(*TableName), nil
}

func (font *Font) HheaTable() (*TableHhea, error) {
	t, err := font.Table(TagHhea)
	if err != nil {
		return nil, err
	}
	return t.(*TableHhea), nil
}

func (font *Font) OS2Table() (*TableOS2, error) {
	t, err := font.Table(TagOS2)
	if err != nil {
		return nil, err
	}
	return t.(*TableOS2), nil
}

func (font *Font) TableLayout(tag Tag) (*TableLayout, error) {
	t, err := font.Table(tag)
	if err != nil {
		return nil, err
	}
	l, ok := t.(*TableLayout)
	if !ok {
		return nil, fmt.Errorf("table %q is not a layout table", tag)
	}
	return l, nil
}

// GposTable returns the Glyph Positioning table identified with the 'GPOS' tag.
func (font *Font) GposTable() (*TableLayout, error) {
	return font.TableLayout(TagGpos)
}

// GsubTable returns the Glyph Substitution table identified with the 'GSUB' tag.
func (font *Font) GsubTable() (*TableLayout, error) {
	return font.TableLayout(TagGsub)
}

func (font *Font) Table(tag Tag) (Table, error) {
	s, found := font.tables[tag]
	if !found {
		return nil, ErrMissingTable
	}

	if s.table == nil {
		t, err := font.parseTable(s)
		if err != nil {
			return nil, err
		}
		s.table = t
	}
	return s.table, nil
}

// New returns an empty Font. It has only an empty 'head' table.
func New(scalerType Tag) *Font {
	font := &Font{
		scalerType: scalerType,
		tables:     make(map[Tag]*tableSection),
	}
	font.AddTable(TagHead, &TableHead{})
	return font
}

// File is a combination of io.Reader, io.Seeker and io.ReaderAt.
// This interface is satisfied by most things that you'd want
// to parse, for example os.File, or io.SectionReader.
type File interface {
	Read([]byte) (int, error)
	ReadAt([]byte, int64) (int, error)
	Seek(int64, int) (int64, error)
}

// Parse parses an OpenType, TrueType, WOFF, or WOFF2 file and returns a Font.
// If parsing fails, an error is returned and *Font will be nil.
func Parse(file File) (*Font, error) {
	magic, err := ReadTag(file)
	if err != nil {
		return nil, err
	}

	file.Seek(0, 0)

	switch magic {
	case SignatureWOFF:
		return parseWOFF(file)
	case SignatureWOFF2:
		return parseWOFF2(file)
	case TypeTrueType, TypeOpenType, TypePostScript1, TypeAppleTrueType:
		return parseOTF(file)
	default:
		return nil, ErrUnsupportedFormat
	}
}

// StrictParse parses an OpenType, TrueType, WOFF or WOFF2 file and returns a Font.
// Each table will be fully parsed and an error is returned if any fail.
func StrictParse(file File) (*Font, error) {
	font, err := Parse(file)
	if err != nil {
		return nil, err
	}

	for _, tag := range font.Tags() {
		if _, err := font.Table(tag); err != nil {
			return nil, fmt.Errorf("failed to parse %q: %s", tag, err)
		}
	}

	return font, nil
}
