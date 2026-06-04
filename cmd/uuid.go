package cmd

import (
	"fmt"

	"github.com/benstro/dt/internal/uuid"
	"github.com/spf13/cobra"
)

var uuidCmd = &cobra.Command{
	Use:   "uuid",
	Short: "Generate a UUID v4",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println(uuid.Generate())
		return nil
	},
}

func init() {
	rootCmd.AddCommand(uuidCmd)
}
