package cmd

import (
	"fmt"
	"os"

	"github.com/aura-studio/dynamic-cli/meta"
	"github.com/spf13/cobra"
)

var metaCallCmd = &cobra.Command{
	Use:   "call <file.so>",
	Short: "Extract meta by calling exported function",
	Long: `Load the .so file as a Go plugin and call the exported Meta() function.

This is the most reliable method as it directly invokes the Meta() function
that returns a JSON string with all metadata. However, it only works with
libgo_*.so files (Go plugins), not libcgo_*.so files (C shared libraries).

The plugin must export a 'Tunnel' variable with a Meta() method that returns
a JSON string containing the metadata fields.

Requirements:
  - The .so file must be a Go plugin (built with -buildmode=plugin)
  - The plugin must be compatible with the current Go version and platform

Example:
  dynamic meta call libgo_xxx.so`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		soPath := args[0]
		if _, err := os.Stat(soPath); os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "file not found: %s\n", soPath)
			os.Exit(1)
		}

		result, jsonStr, err := meta.ReadCall(soPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to read meta: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Raw JSON:")
		fmt.Println(jsonStr)
		fmt.Println()
		fmt.Println("Parsed:")
		for _, key := range meta.Keys {
			if val, ok := result[key]; ok {
				fmt.Printf("  %s: %s\n", key, val)
			}
		}
	},
}

func init() {
	metaCmd.AddCommand(metaCallCmd)
}
