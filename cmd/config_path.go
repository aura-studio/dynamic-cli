package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// resolveConfigPath returns the config path to use.
// Priority:
// 1) explicit --config/-c
// 2) ./dynamic-cli.yaml (if it exists in current directory)
// 3) ./dynamic-cli.yml (if it exists in current directory)
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

	// Try new default name first
	const defaultCliYaml = "dynamic-cli.yaml"
	if _, err := os.Stat(defaultCliYaml); err == nil {
		return defaultCliYaml
	}
	const defaultCliYml = "dynamic-cli.yml"
	if _, err := os.Stat(defaultCliYml); err == nil {
		return defaultCliYml
	}

	fmt.Println("error: config file not found (use -c <path>), or place dynamic-cli.yaml/dynamic-cli.yml in current directory")
	os.Exit(1)
	return ""
}
