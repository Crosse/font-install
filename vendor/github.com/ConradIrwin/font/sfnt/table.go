package sfnt

import (
	"compress/zlib"
	"io"
)

var parsers = map[Tag]tableParser{
	TagHead: parseTableHead,
	TagName: parseTableName,
	TagHhea: parseTableHhea,
	TagOS2:  parseTableOS2,
	TagGpos: parseTableLayout,
	TagGsub: parseTableLayout,
}

// Table is an interface for each section of the font file.
type Table interface {
	Bytes() []byte
	Name() string // Name returns the name of the table.
}

type baseTable Tag

// Name returns the name of the table represented by this tag.
func (b baseTable) Name() string {
	return tableTags[Tag(b).String()]
}

type unparsedTable struct {
	baseTable

	bytes []byte // Uncompress content of this table.
}

type tableParser func(tag Tag, buffer []byte) (Table, error)

func newUnparsedTable(tag Tag, buffer []byte) (Table, error) {
	return &unparsedTable{baseTable(tag), buffer}, nil
}

func (font *Font) parseTable(s *tableSection) (Table, error) {
	var buf []byte

	if s.length != 0 && s.length < s.zLength {
		zbuf := io.NewSectionReader(font.file, int64(s.offset), int64(s.length))
		r, err := zlib.NewReader(zbuf)
		if err != nil {
			return nil, err
		}
		defer r.Close()

		buf = make([]byte, s.zLength, s.zLength)
		if _, err := io.ReadFull(r, buf); err != nil {
			return nil, err
		}
	} else {
		buf = make([]byte, s.length, s.length)
		if _, err := font.file.ReadAt(buf, int64(s.offset)); err != nil {
			return nil, err
		}
	}

	parser, found := parsers[s.tag]
	if !found {
		parser = newUnparsedTable
	}

	return parser(s.tag, buf)
}
