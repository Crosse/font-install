package sfnt

import (
	"bytes"
	"testing"
)

func TestParseCrashers(t *testing.T) {

	font, err := Parse(bytes.NewReader([]byte{}))
	if font != nil || err == nil {
		t.Fail()
	}

}
