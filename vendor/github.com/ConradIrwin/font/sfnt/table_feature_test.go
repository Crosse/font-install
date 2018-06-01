package sfnt

import (
	"testing"
)

func TestFeatureString(t *testing.T) {
	tests := []struct {
		tag  string
		want string
	}{
		{"aalt", "Access All Alternates"},
		{"zhtw", "Traditional Chinese Forms (Deprecated)"},
		{"ss01", "Stylistic Set 1"},
		{"ss20", "Stylistic Set 20"},
		{"cv01", "Character Variant 1"},
		{"cv99", "Character Variant 99"},
		{"BLAH", ""},
		{"cv00", ""},
		{"ss00", ""},
		{"ss21", ""},
	}

	for _, test := range tests {
		f := Feature{Tag: MustNamedTag(test.tag)}
		if got := f.String(); got != test.want {
			t.Errorf("Feature{%q}.String() = %q want %q", test.tag, got, test.want)
		}
	}
}
