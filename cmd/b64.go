package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/benstro/dt/internal/b64"
	"github.com/spf13/cobra"
)

var b64Cmd = &cobra.Command{
	Use:   "b64",
	Short: "Base64 encode and decode",
}

var b64EncodeCmd = &cobra.Command{
	Use:   "encode [text]",
	Short: "Encode text to base64",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		input, err := resolveInput(args)
		if err != nil {
			return err
		}
		fmt.Println(b64.Encode(input))
		return nil
	},
}

var b64DecodeCmd = &cobra.Command{
	Use:   "decode [text]",
	Short: "Decode base64 text",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		input, err := resolveInput(args)
		if err != nil {
			return err
		}
		result, err := b64.Decode(input)
		if err != nil {
			return fmt.Errorf("decode failed: %w", err)
		}
		fmt.Println(result)
		return nil
	},
}

func resolveInput(args []string) (string, error) {
	if len(args) == 1 {
		return args[0], nil
	}
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			return "", fmt.Errorf("reading stdin: %w", err)
		}
		return string(data), nil
	}
	return "", fmt.Errorf("provide input as an argument or via stdin")
}

func init() {
	b64Cmd.AddCommand(b64EncodeCmd, b64DecodeCmd)
	rootCmd.AddCommand(b64Cmd)
}
