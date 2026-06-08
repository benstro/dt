package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func resolveInput(args []string) (string, error) {
	return resolveInputRaw(args, false)
}

func resolveRawInput(args []string) (string, error) {
	return resolveInputRaw(args, true)
}

func resolveInputRaw(args []string, raw bool) (string, error) {
	if len(args) == 1 {
		return args[0], nil
	}

	stat, _ := os.Stdin.Stat()
	if stat.Mode()&os.ModeCharDevice == 0 {
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			return "", fmt.Errorf("reading stdin: %w", err)
		}
		if raw {
			return string(data), nil
		}
		return strings.TrimRight(string(data), "\n"), nil
	}
	return "", fmt.Errorf("provide input as an argument or via stdin")
}
