// Package json provides JSON formatting utilities.
package json

import (
	stdjson "encoding/json"
	"fmt"
)

func Pretty(input string) (string, error) {
	var v any
	if err := stdjson.Unmarshal([]byte(input), &v); err != nil {
		return "", fmt.Errorf("invalid JSON: %w", err)
	}
	b, err := stdjson.MarshalIndent(v, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func Minify(input string) (string, error) {
	var v any
	if err := stdjson.Unmarshal([]byte(input), &v); err != nil {
		return "", fmt.Errorf("invalid JSON: %w", err)
	}
	b, err := stdjson.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func Stringify(input string) (string, error) {
	var v any
	if err := stdjson.Unmarshal([]byte(input), &v); err != nil {
		return "", fmt.Errorf("invalid JSON: %w", err)
	}
	b, err := stdjson.Marshal(v)
	if err != nil {
		return "", err
	}
	escaped, err := stdjson.Marshal(string(b))
	if err != nil {
		return "", err
	}
	return string(escaped), nil
}
