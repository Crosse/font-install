package sfnt

import (
	"encoding/binary"
	"io"
	"sort"
)

// these headers always seem to come first in the serialized output.
var outputOrder = map[Tag]int{
	TagMaxp: 0,
	TagHead: 1,
	TagHmtx: 2,
	TagHhea: 3,
	TagOS2:  4,
	TagName: 5,
}

// WriteOTF serializes a Font into OpenType format suitable
// for writing to a file such as *.otf.
// You can also use this to write to files called *.ttf if the
// font contains TrueType glyphs.
func (font *Font) WriteOTF(w io.Writer) (n int, err error) {

	todo := font.Tags()
	sort.Slice(todo, func(i, j int) bool {
		iScore, ok := outputOrder[todo[i]]
		if !ok {
			iScore = int(todo[i].Number)
		}
		jScore, ok := outputOrder[todo[j]]
		if !ok {
			jScore = int(todo[j].Number)
		}

		return iScore < jScore
	})

	headTable, err := font.HeadTable()
	if err != nil {
		return n, err
	}

	headTable.ClearExpectedChecksum()

	header := newOTFHeader(font.scalerType, uint16(len(todo)))

	fragments := make([][]byte, len(todo))

	offset := otfHeaderLength + directoryEntryLength*len(todo)
	checksum := header.checkSum()

	err = binary.Write(w, binary.BigEndian, header)
	if err != nil {
		return n, err
	}
	n += otfHeaderLength

	for i, tag := range todo {
		t, err := font.Table(tag)
		if err != nil {
			return n, err
		}
		fragments[i] = t.Bytes()
		entry := directoryEntry{
			Tag:      tag,
			CheckSum: checkSum(fragments[i]),
			Offset:   uint32(offset),
			Length:   uint32(len(fragments[i])),
		}

		offset += len(fragments[i])
		if len(fragments[i])%4 != 0 {
			offset += 4 - (len(fragments[i]) % 4)

		}
		checksum += entry.CheckSum + entry.checkSum()

		err = binary.Write(w, binary.BigEndian, entry)
		if err != nil {
			return n, err
		}
		n += directoryEntryLength
	}

	for i, tag := range todo {

		var fragment []byte

		if tag == TagHead {
			headTable.SetExpectedChecksum(checksum)
			fragment = headTable.Bytes()
			headTable.SetExpectedChecksum(0)
		} else {
			fragment = fragments[i]
		}

		m, err := w.Write(fragment)
		n += m
		if err != nil {
			return n, err
		}

		var extra []byte

		if len(fragment)%4 == 1 {
			extra = []byte{0, 0, 0}

		} else if len(fragment)%4 == 2 {
			extra = []byte{0, 0}
		} else if len(fragment)%4 == 3 {
			extra = []byte{0}
		} else {
			continue
		}

		m, err = w.Write(extra)
		n += m
		if err != nil {
			return n, err
		}

	}

	return 0, nil
}

func checkSum(buffer []byte) uint32 {
	total := uint32(0)

	for len(buffer) >= 4 {
		total += uint32(buffer[0])<<24 | uint32(buffer[1])<<16 | uint32(buffer[2])<<8 | uint32(buffer[3])
		buffer = buffer[4:]
	}

	if len(buffer) >= 1 {
		total = total + uint32(buffer[0])<<24
	}
	if len(buffer) >= 2 {
		total = total + uint32(buffer[1])<<16
	}
	if len(buffer) >= 3 {
		total += uint32(buffer[2]) << 8
	}

	return total

}
