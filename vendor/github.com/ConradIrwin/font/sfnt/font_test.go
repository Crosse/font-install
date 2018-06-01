package sfnt

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

// TestSmokeTest very simply checks we can parse, and write the sample fonts
// without error.
// TODO We should check what is returned is valid.
func TestSmokeTest(t *testing.T) {
	tests := []struct {
		filename string
	}{
		{filename: "Roboto-BoldItalic.ttf"},
		{filename: "Raleway-v4020-Regular.otf"},
		{filename: "open-sans-v15-latin-regular.woff"},
		{filename: "Go-Regular.woff2"},
	}

	for _, test := range tests {
		filename := filepath.Join("testdata", test.filename)
		file, err := os.Open(filename)
		if err != nil {
			t.Errorf("Failed to open %q: %s\n", filename, err)
		}

		font, err := StrictParse(file)
		if err != nil {
			t.Errorf("StrictParse(%q) err = %q, want nil", filename, err)
			continue
		}

		if _, err := font.WriteOTF(ioutil.Discard); err != nil {
			t.Errorf("WriteOTF(%q) err = %q, want nil", filename, err)
			continue
		}

		file.Close()
	}
}

// benchmarkParse tests the performance of a simple Parse.
// Example run:
//   go test -cpuprofile cpu.prof -benchmem -memprofile mem.prof -bench . -run=^$ -benchtime=30s github.com/ConradIrwin/font/sfnt
//   go tool pprof cpu.prof
//
// BenchmarkParseOTF-8           	20000000	      3209 ns/op	    1229 B/op	      32 allocs/op
// BenchmarkStrictParseOTF-8     	  200000	    184822 ns/op	  372415 B/op	    1616 allocs/op
// BenchmarkParseWOFF-8          	10000000	      3999 ns/op	    1993 B/op	      40 allocs/op
// BenchmarkStrictParseWOFF-8    	   50000	    776500 ns/op	  575990 B/op	     497 allocs/op
// BenchmarkParseWOFF2-8         	   20000	   2011769 ns/op	  742531 B/op	     468 allocs/op
// BenchmarkStrictParseWOFF2-8   	   20000	   2033596 ns/op	  875608 B/op	     818 allocs/op
func benchmarkParse(b *testing.B, filename string) {
	buf, err := ioutil.ReadFile(filepath.Join("testdata", filename))
	if err != nil {
		b.Errorf("Failed to open %q: %s\n", filename, err)
	}

	for n := 0; n < b.N; n++ {
		r := bytes.NewReader(buf)
		if _, err := Parse(r); err != nil {
			b.Errorf("Parse(%q) err = %q, want nil", filename, err)
			return
		}
	}
}

// benchmarkStrictParse tests the performance of a simple StrictParse.
func benchmarkStrictParse(b *testing.B, filename string) {
	buf, err := ioutil.ReadFile(filepath.Join("testdata", filename))
	if err != nil {
		b.Errorf("Failed to open %q: %s\n", filename, err)
	}

	for n := 0; n < b.N; n++ {
		r := bytes.NewReader(buf)
		if _, err := StrictParse(r); err != nil {
			b.Errorf("StrictParse(%q) err = %q, want nil", filename, err)
			return
		}
	}
}

func BenchmarkParseOTF(b *testing.B) {
	benchmarkParse(b, "Roboto-BoldItalic.ttf")
}

func BenchmarkStrictParseOTF(b *testing.B) {
	benchmarkStrictParse(b, "Roboto-BoldItalic.ttf")
}

func BenchmarkParseWOFF(b *testing.B) {
	benchmarkParse(b, "open-sans-v15-latin-regular.woff")
}

func BenchmarkStrictParseWOFF(b *testing.B) {
	benchmarkStrictParse(b, "open-sans-v15-latin-regular.woff")
}

func BenchmarkParseWOFF2(b *testing.B) {
	benchmarkParse(b, "Go-Regular.woff2")
}

func BenchmarkStrictParseWOFF2(b *testing.B) {
	benchmarkStrictParse(b, "Go-Regular.woff2")
}
