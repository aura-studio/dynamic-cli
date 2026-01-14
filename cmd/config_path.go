package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// resolveConfigPath returns the config path to use.
// Priority:
// 1) explicit --config/-c
// 2) ./dynamic.yaml (if it exists in current directory)
// 3) ./dynamic.yml (if it exists in current directory)
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

	// Try dynamic.yaml first
	const defaultYaml = "dynamic.yaml"
	if _, err := os.Stat(defaultYaml); err == nil {
		return defaultYaml
	}

	// Try dynamic.yml as fallback
	const defaultYml = "dynamic.yml"
	if _, err := os.Stat(defaultYml); err == nil {
		return defaultYml
	}

	fmt.Println("error: config file not found (use -c <path>), or place dynamic.yaml/dynamic.yml in current directory")
	os.Exit(1)
	return ""
}
