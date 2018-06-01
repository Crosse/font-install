package woff2

import "testing"

func TestKnownTableTagsLength(t *testing.T) {
	const want = 63
	if got := len(knownTableTags); got != want {
		t.Errorf("got len(knownTableTags): %v, want: %v", got, want)
	}
}
