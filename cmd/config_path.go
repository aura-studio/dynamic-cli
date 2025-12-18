package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// resolveConfigPath returns the config path to use.
// Priority:
// 1) explicit --config/-c
// 2) ./dynamic.yaml (only if it exists in current directory)
// Otherwise it prints an error and exits with code 1.
func resolveConfigPath(cmd *cobra.Command) string {
	cfgPath, err := cmd.Flags().GetString("config")
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	if cfgPath != "" {
		return cfgPath
	}

	const defaultName = "dynamic.yaml"
	if _, err := os.Stat(defaultName); err == nil {
		return defaultName
	}

	fmt.Println("error: config is required (use -c <path>), or place dynamic.yaml in current directory")
	os.Exit(1)
	return ""
}
