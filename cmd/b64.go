package cmd

import (
	"fmt"

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

func init() {
	b64Cmd.AddCommand(b64EncodeCmd, b64DecodeCmd)
	rootCmd.AddCommand(b64Cmd)
}
