package hash

import (
	"errors"
	"io"
	"os"
	"sort"
	"strings"
	"testing"
)

// errReader is an io.Reader that always returns an error, used to exercise
// the error path in HashSha256File.
type errReader struct{ err error }

func (e errReader) Read(_ []byte) (int, error) { return 0, e.err }

// resultsByPath sorts a []Result slice in-place by Path so that tests on
// HashSha256Files can make deterministic assertions regardless of goroutine
// scheduling order.
func resultsByPath(rs []Result) {
	sort.Slice(rs, func(i, j int) bool { return rs[i].Path < rs[j].Path })
}

// writeTempFile creates a named temp file inside dir, writes content to it,
// and returns the file path.
func writeTempFile(t *testing.T, dir, content string) string {
	t.Helper()
	f, err := os.CreateTemp(dir, "hash-test-*")
	if err != nil {
		t.Fatalf("os.CreateTemp: %v", err)
	}
	if _, err := io.WriteString(f, content); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	if err := f.Close(); err != nil {
		t.Fatalf("close temp file: %v", err)
	}
	return f.Name()
}

// ---------------------------------------------------------------------------
// HashSha256
// ---------------------------------------------------------------------------

func TestHashSha256_KnownVectors(t *testing.T) {
	// SHA-256 reference values from NIST / Wikipedia.
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "empty string",
			input: "",
			want:  "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
		},
		{
			name:  "hello",
			input: "hello",
			want:  "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824",
		},
		{
			name:  "The quick brown fox jumps over the lazy dog",
			input: "The quick brown fox jumps over the lazy dog",
			want:  "d7a8fbb307d7809469ca9abcb0082e4f8d5651e46d3cdb762d02d0bf37c9e592",
		},
		{
			name:  "single space",
			input: " ",
			want:  "36a9e7f1c95b82ffb99743e0c5c4ce95d83c9a430aac59f84ef3cbfab6145068",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := HashSha256(tc.input)
			if got != tc.want {
				t.Errorf("HashSha256(%q) = %q, want %q", tc.input, got, tc.want)
			}
		})
	}
}

func TestHashSha256_OutputLength(t *testing.T) {
	// SHA-256 hex digest is always 64 characters.
	got := HashSha256("anything")
	if len(got) != 64 {
		t.Errorf("HashSha256 output length = %d, want 64", len(got))
	}
}

func TestHashSha256_IsDeterministic(t *testing.T) {
	input := "determinism check"
	if HashSha256(input) != HashSha256(input) {
		t.Error("HashSha256 returned different results for identical input")
	}
}

// ---------------------------------------------------------------------------
// HashSha256File
// ---------------------------------------------------------------------------

func TestHashSha256File_KnownContent(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    string
	}{
		{
			name:    "empty reader",
			content: "",
			want:    "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
		},
		{
			name:    "hello",
			content: "hello",
			want:    "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824",
		},
		{
			name:    "multiline content",
			content: "line1\nline2\n",
			want:    "2751a3a2f303ad21752038085e2b8c5f98ecff61a2e4ebbd43506a941725be80",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := HashSha256File(strings.NewReader(tc.content))
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tc.want {
				t.Errorf("HashSha256File(%q) = %q, want %q", tc.content, got, tc.want)
			}
		})
	}
}

func TestHashSha256File_MatchesHashSha256(t *testing.T) {
	// HashSha256File and HashSha256 must agree on the same input.
	input := "consistency check"
	fromString := HashSha256(input)
	fromReader, err := HashSha256File(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fromString != fromReader {
		t.Errorf("HashSha256(%q)=%q differs from HashSha256File(%q)=%q",
			input, fromString, input, fromReader)
	}
}

func TestHashSha256File_FailingReader(t *testing.T) {
	sentinel := errors.New("disk exploded")
	_, err := HashSha256File(errReader{err: sentinel})
	if err == nil {
		t.Fatal("expected an error from a failing reader, got nil")
	}
	if !errors.Is(err, sentinel) {
		t.Errorf("error chain does not contain sentinel: %v", err)
	}
}

// ---------------------------------------------------------------------------
// HashSha256Files
// ---------------------------------------------------------------------------

func TestHashSha256Files_EmptySlice(t *testing.T) {
	results := HashSha256Files([]string{})
	if len(results) != 0 {
		t.Errorf("expected 0 results for empty input, got %d", len(results))
	}
}

func TestHashSha256Files_SingleFile(t *testing.T) {
	dir := t.TempDir()
	content := "hello"
	path := writeTempFile(t, dir, content)

	results := HashSha256Files([]string{path})

	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	r := results[0]
	if r.Err != nil {
		t.Fatalf("unexpected error for %q: %v", path, r.Err)
	}
	want := HashSha256(content)
	if r.Hash != want {
		t.Errorf("HashSha256Files single file: got hash %q, want %q", r.Hash, want)
	}
	if r.Path != path {
		t.Errorf("HashSha256Files single file: got path %q, want %q", r.Path, path)
	}
}

func TestHashSha256Files_MultipleFiles(t *testing.T) {
	dir := t.TempDir()

	files := []struct {
		content string
	}{
		{"alpha content"},
		{"beta content"},
		{"gamma content"},
	}

	var paths []string
	wantHashes := make(map[string]string)
	for _, f := range files {
		p := writeTempFile(t, dir, f.content)
		paths = append(paths, p)
		wantHashes[p] = HashSha256(f.content)
	}

	results := HashSha256Files(paths)

	if len(results) != len(paths) {
		t.Fatalf("expected %d results, got %d", len(paths), len(results))
	}

	resultsByPath(results)

	for _, r := range results {
		if r.Err != nil {
			t.Errorf("unexpected error for %q: %v", r.Path, r.Err)
			continue
		}
		want, ok := wantHashes[r.Path]
		if !ok {
			t.Errorf("result for unexpected path %q", r.Path)
			continue
		}
		if r.Hash != want {
			t.Errorf("path %q: got hash %q, want %q", r.Path, r.Hash, want)
		}
	}
}

func TestHashSha256Files_NonExistentPath(t *testing.T) {
	missing := "/tmp/this-file-does-not-exist-hash-test-xyz123"

	results := HashSha256Files([]string{missing})

	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	r := results[0]
	if r.Path != missing {
		t.Errorf("got path %q, want %q", r.Path, missing)
	}
	if r.Err == nil {
		t.Error("expected an error for non-existent file, got nil")
	}
	if !errors.Is(r.Err, os.ErrNotExist) {
		t.Errorf("expected os.ErrNotExist in error chain, got: %v", r.Err)
	}
	if r.Hash != "" {
		t.Errorf("expected empty hash on error, got %q", r.Hash)
	}
}

func TestHashSha256Files_MixedExistentAndMissing(t *testing.T) {
	dir := t.TempDir()
	existingPath := writeTempFile(t, dir, "real content")
	missingPath := dir + "/does-not-exist"

	results := HashSha256Files([]string{existingPath, missingPath})

	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}

	byPath := make(map[string]Result, 2)
	for _, r := range results {
		byPath[r.Path] = r
	}

	good, ok := byPath[existingPath]
	if !ok {
		t.Fatalf("no result for existing path %q", existingPath)
	}
	if good.Err != nil {
		t.Errorf("unexpected error for existing file: %v", good.Err)
	}
	if good.Hash != HashSha256("real content") {
		t.Errorf("wrong hash for existing file: %q", good.Hash)
	}

	bad, ok := byPath[missingPath]
	if !ok {
		t.Fatalf("no result for missing path %q", missingPath)
	}
	if bad.Err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

func TestHashSha256Files_ResultCountMatchesInput(t *testing.T) {
	// Regardless of success/failure, every input path must produce exactly
	// one result — verifying the fan-out/fan-in contract.
	dir := t.TempDir()
	paths := []string{
		writeTempFile(t, dir, "a"),
		writeTempFile(t, dir, "b"),
		"/tmp/missing-hash-test-1",
		writeTempFile(t, dir, "c"),
		"/tmp/missing-hash-test-2",
	}

	results := HashSha256Files(paths)

	if len(results) != len(paths) {
		t.Errorf("got %d results for %d inputs", len(results), len(paths))
	}
}
