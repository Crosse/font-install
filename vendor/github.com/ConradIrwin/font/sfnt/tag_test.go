package sfnt

import "testing"

func TestParsedTag(t *testing.T) {
	tag := MustNamedTag("head")
	if tag.Number != 0x68656164 {
		t.Errorf("head != %v", tag.Number)
	}

	if tag.String() != "head" {
		t.Errorf("head != %v", tag.String())
	}
}

func TestNewTag(t *testing.T) {
	tag := Tag{0x74727565}

	if tag.Number != 0x74727565 {
		t.Errorf("true != %v", tag.Number)
	}

	if tag.String() != "true" {
		t.Errorf("true != %v", tag.Number)
	}
}

func TestTagEquality(t *testing.T) {
	t1 := Tag{0x74727565}
	t2 := Tag{0x74727565}

	if t1 != t2 {
		t.Errorf("equality failed for number")
	}

	if MustNamedTag("head") != MustNamedTag("head") {
		t.Errorf("equality failed for parsed")
	}

	if MustNamedTag("true") != t1 {
		t.Errorf("equality failed %v %v", MustNamedTag("true"), t1)
	}
}
