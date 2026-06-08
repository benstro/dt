// Package hash provides hashing utilities
package hash

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"sync"
)

type Result struct {
	Path string
	Hash string
	Err  error
}

func HashSha256(data string) string {
	hash := sha256.Sum256([]byte(data))
	hexHash := hex.EncodeToString(hash[:])
	return hexHash
}

func HashSha256Files(paths []string) []Result {
	results := make(chan Result, len(paths))
	var wg sync.WaitGroup

	for _, p := range paths {
		wg.Add(1)
		go func(path string) {
			defer wg.Done()
			f, err := os.Open(path)
			if err != nil {
				results <- Result{Path: path, Err: err}
				return
			}
			defer f.Close()
			h, err := HashSha256File(f)
			results <- Result{Path: path, Hash: h, Err: err}
		}(p)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	var out []Result
	for r := range results {
		out = append(out, r)
	}
	return out
}

func HashSha256File(r io.Reader) (string, error) {
	h := sha256.New()
	_, err := io.Copy(h, r)
	if err != nil {
		return "", fmt.Errorf("error copying from file: %w", err)
	}
	sum := h.Sum(nil)

	return hex.EncodeToString(sum), nil
}
