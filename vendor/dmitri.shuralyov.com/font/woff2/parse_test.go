package woff2_test

import (
	"fmt"
	"log"

	"dmitri.shuralyov.com/font/woff2"
	"github.com/shurcooL/gofontwoff"
)

func ExampleParse() {
	f, err := gofontwoff.Assets.Open("/Go-Regular.woff2")
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()

	font, err := woff2.Parse(f)
	if err != nil {
		log.Fatalln(err)
	}
	Dump(font)

	// Output:
	//
	// Signature:           0x774f4632
	// Flavor:              0x00010000
	// Length:              46132
	// NumTables:           14
	// Reserved:            0
	// TotalSfntSize:       140308
	// TotalCompressedSize: 46040
	// MajorVersion:        1
	// MinorVersion:        0
	// MetaOffset:          0
	// MetaLength:          0
	// MetaOrigLength:      0
	// PrivOffset:          0
	// PrivLength:          0
	//
	// TableDirectory: 14 entries
	// 	{Flags: 0x06, Tag: <nil>, OrigLength: 96, TransformLength: <nil>}
	// 	{Flags: 0x00, Tag: <nil>, OrigLength: 1318, TransformLength: <nil>}
	// 	{Flags: 0x08, Tag: <nil>, OrigLength: 176, TransformLength: <nil>}
	// 	{Flags: 0x09, Tag: <nil>, OrigLength: 3437, TransformLength: <nil>}
	// 	{Flags: 0x11, Tag: <nil>, OrigLength: 8, TransformLength: <nil>}
	// 	{Flags: 0x0a, Tag: <nil>, OrigLength: 118912, TransformLength: 105020}
	// 	{Flags: 0x0b, Tag: <nil>, OrigLength: 1334, TransformLength: 0}
	// 	{Flags: 0x01, Tag: <nil>, OrigLength: 54, TransformLength: <nil>}
	// 	{Flags: 0x02, Tag: <nil>, OrigLength: 36, TransformLength: <nil>}
	// 	{Flags: 0x03, Tag: <nil>, OrigLength: 2662, TransformLength: <nil>}
	// 	{Flags: 0x04, Tag: <nil>, OrigLength: 32, TransformLength: <nil>}
	// 	{Flags: 0x05, Tag: <nil>, OrigLength: 6967, TransformLength: <nil>}
	// 	{Flags: 0x07, Tag: <nil>, OrigLength: 4838, TransformLength: <nil>}
	// 	{Flags: 0x0c, Tag: <nil>, OrigLength: 188, TransformLength: <nil>}
	//
	// CollectionDirectory: <nil>
	// CompressedFontData: 124832 bytes (uncompressed size)
	// ExtendedMetadata: <nil>
	// PrivateData: []
}

func Dump(f woff2.File) {
	dumpHeader(f.Header)
	fmt.Println()
	dumpTableDirectory(f.TableDirectory)
	fmt.Println()
	fmt.Println("CollectionDirectory:", f.CollectionDirectory)
	fmt.Println("CompressedFontData:", len(f.FontData), "bytes (uncompressed size)")
	fmt.Println("ExtendedMetadata:", f.ExtendedMetadata)
	fmt.Println("PrivateData:", f.PrivateData)
}

func dumpHeader(hdr woff2.Header) {
	fmt.Printf("Signature:           %#08x\n", hdr.Signature)
	fmt.Printf("Flavor:              %#08x\n", hdr.Flavor)
	fmt.Printf("Length:              %d\n", hdr.Length)
	fmt.Printf("NumTables:           %d\n", hdr.NumTables)
	fmt.Printf("Reserved:            %d\n", hdr.Reserved)
	fmt.Printf("TotalSfntSize:       %d\n", hdr.TotalSfntSize)
	fmt.Printf("TotalCompressedSize: %d\n", hdr.TotalCompressedSize)
	fmt.Printf("MajorVersion:        %d\n", hdr.MajorVersion)
	fmt.Printf("MinorVersion:        %d\n", hdr.MinorVersion)
	fmt.Printf("MetaOffset:          %d\n", hdr.MetaOffset)
	fmt.Printf("MetaLength:          %d\n", hdr.MetaLength)
	fmt.Printf("MetaOrigLength:      %d\n", hdr.MetaOrigLength)
	fmt.Printf("PrivOffset:          %d\n", hdr.PrivOffset)
	fmt.Printf("PrivLength:          %d\n", hdr.PrivLength)
}

func dumpTableDirectory(td woff2.TableDirectory) {
	fmt.Println("TableDirectory:", len(td), "entries")
	for _, t := range td {
		fmt.Printf("\t{")
		fmt.Printf("Flags: %#02x, ", t.Flags)
		if t.Tag != nil {
			fmt.Printf("Tag: %v, ", *t.Tag)
		} else {
			fmt.Printf("Tag: <nil>, ")
		}
		fmt.Printf("OrigLength: %v, ", t.OrigLength)
		if t.TransformLength != nil {
			fmt.Printf("TransformLength: %v", *t.TransformLength)
		} else {
			fmt.Printf("TransformLength: <nil>")
		}
		fmt.Printf("}\n")
	}
}
