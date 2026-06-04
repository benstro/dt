package cmd

import (
	"fmt"

	dtjson "github.com/benstro/dt/internal/json"
	"github.com/spf13/cobra"
)

var jsonCmd = &cobra.Command{
	Use:   "json",
	Short: "JSON manipulation tools",
}

var jsonPrettyCmd = &cobra.Command{
	Use:   "pretty <json>",
	Short: "Format JSON with indentation",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		result, err := dtjson.Pretty(args[0])
		if err != nil {
			return err
		}
		fmt.Println(result)
		return nil
	},
}

var jsonMinifyCmd = &cobra.Command{
	Use:   "minify <json>",
	Short: "Strip whitespace from JSON",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		result, err := dtjson.Minify(args[0])
		if err != nil {
			return err
		}
		fmt.Println(result)
		return nil
	},
}

var jsonStringifyCmd = &cobra.Command{
	Use:   "stringify <json>",
	Short: "Escape a JSON value into a JSON string literal",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		result, err := dtjson.Stringify(args[0])
		if err != nil {
			return err
		}
		fmt.Println(result)
		return nil
	},
}

func init() {
	jsonCmd.AddCommand(jsonPrettyCmd, jsonMinifyCmd, jsonStringifyCmd)
	rootCmd.AddCommand(jsonCmd)
}
