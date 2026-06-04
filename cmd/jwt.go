package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/benstro/dt/internal/jwt"
	"github.com/spf13/cobra"
)

var jwtCmd = &cobra.Command{
	Use:   "jwt",
	Short: "Inspect JSON Web Tokens",
}

var jwtDecodeCmd = &cobra.Command{
	Use:   "decode <token>",
	Short: "Decode and pretty-print a JWT header and payload",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		token, err := jwt.Decode(args[0])
		if err != nil {
			return err
		}

		header, _ := json.MarshalIndent(token.Header, "", "  ")
		payload, _ := json.MarshalIndent(token.Payload, "", "  ")

		fmt.Printf("Header:\n%s\n\nPayload:\n%s\n", header, payload)
		return nil
	},
}

func init() {
	jwtCmd.AddCommand(jwtDecodeCmd)
	rootCmd.AddCommand(jwtCmd)
}
