package cmd

import (
	"fmt"
	"os"

	"github.com/benstro/dt/internal/hash"
	"github.com/spf13/cobra"
)

var hashCmd = &cobra.Command{
	Use:   "hash",
	Short: "Hash a string or file",
}

var filePath string

var hashSha256Cmd = &cobra.Command{
	Use:   "sha256 [string | file...]",
	Short: "Hash to sha256",
	RunE: func(cmd *cobra.Command, args []string) error {
		// --file flag or multiple positional args → parallel file hashing
		if filePath != "" {
			return hashFiles([]string{filePath})
		}
		if len(args) > 1 {
			return hashFiles(args)
		}
		return hashFromString(args)
	},
}

func hashFiles(paths []string) error {
	results := hash.HashSha256Files(paths)
	for _, r := range results {
		if r.Err != nil {
			fmt.Fprintf(os.Stderr, "error: %s: %v\n", r.Path, r.Err)
			continue
		}
		fmt.Printf("%s  %s\n", r.Hash, r.Path)
	}
	return nil
}

func hashFromString(args []string) error {
	input, err := resolveRawInput(args)
	if err != nil {
		return err
	}

	fmt.Println(hash.HashSha256(input))
	return nil
}

func init() {
	hashSha256Cmd.Flags().StringVarP(&filePath, "file", "f", "", "Path for file to hash")
	hashCmd.AddCommand(hashSha256Cmd)
	rootCmd.AddCommand(hashCmd)
}
