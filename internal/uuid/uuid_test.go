package uuid

import (
	"regexp"
	"testing"
)

// uuidV4RE matches the canonical UUID v4 string format:
// xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx
// where y is one of 8, 9, a, or b.
var uuidV4RE = regexp.MustCompile(
	`^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`,
)

func TestGenerate_ReturnsNonEmptyString(t *testing.T) {
	got := Generate()
	if got == "" {
		t.Fatal("Generate() returned an empty string")
	}
}

func TestGenerate_ReturnsValidUUIDv4Format(t *testing.T) {
	got := Generate()
	if !uuidV4RE.MatchString(got) {
		t.Errorf("Generate() = %q, does not match UUID v4 pattern", got)
	}
}

func TestGenerate_TwoCallsReturnDistinctValues(t *testing.T) {
	a := Generate()
	b := Generate()
	if a == b {
		t.Errorf("Generate() returned identical values on consecutive calls: %q", a)
	}
}
