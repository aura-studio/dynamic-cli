/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// cleanCmd represents the clean command
var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean build artifacts",
	Long:  "Cleans build artifacts under the warehouse. Subcommands: cache/package/all/useless. Use -c to specify config; if omitted, ./dynamic-cli.yaml is used only when it exists in the current directory.",
}

func init() {
	rootCmd.AddCommand(cleanCmd)
	cleanCmd.PersistentFlags().StringP("config", "c", "", "path to dynamic-cli.yaml (default: ./dynamic-cli.yaml or ./dynamic-cli.yml)")
	cleanCmd.PersistentFlags().StringP("procedure", "p", "", "procedure name to select warehouse (optional)")
}
