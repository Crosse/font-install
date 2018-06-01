package sfnt

import (
	"bytes"

	"dmitri.shuralyov.com/font/woff2"
)

func parseWOFF2(file File) (*Font, error) {
	f, err := woff2.Parse(file)
	if err != nil {
		return nil, err
	}
	font := &Font{
		file:       bytes.NewReader(f.FontData),
		scalerType: Tag{f.Header.Flavor},
		tables:     make(map[Tag]*tableSection, f.Header.NumTables),
	}
	for _, t := range f.TableDirectory.Tables() {
		tag := Tag{t.Tag}
		font.tables[tag] = &tableSection{
			tag:     tag,
			offset:  uint32(t.Offset),
			length:  uint32(t.Length),
			zLength: uint32(t.Length),
		}
	}
	return font, nil
}
