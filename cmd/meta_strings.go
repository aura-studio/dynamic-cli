package cmd

import (
	"fmt"
	"os"

	"github.com/aura-studio/dynamic-cli/meta"
	"github.com/spf13/cobra"
)

var metaStringsCmd = &cobra.Command{
	Use:   "strings <file.so>",
	Short: "Extract meta using strings command",
	Long: `Use the 'strings' command to search for metadata patterns in the .so file.

This method extracts printable strings from the binary and matches them against
known patterns for each metadata field. It's fast but relies on heuristics,
so results may include false positives or miss values with unexpected formats.

Patterns matched:
  module   - Go module path (e.g. github.com/org/repo)
  version  - Semantic version (e.g. v1.2.3 or pseudo-version)
  built    - Build timestamp (e.g. 2024-06-12_15:30:00_CST+0800)
  os       - Operating system (linux, darwin, windows)
  arch     - Architecture (amd64, arm64, 386)
  compiler - Compiler used (gcc, clang, musl-gcc)
  variant  - Build variant (generic, alpine, debian)

Example:
  dynamic meta strings libcgo_xxx.so`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		soPath := args[0]
		if _, err := os.Stat(soPath); os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "file not found: %s\n", soPath)
			os.Exit(1)
		}

		result, err := meta.ReadStrings(soPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to read meta: %v\n", err)
			os.Exit(1)
		}

		printMeta(result)
	},
}

func init() {
	metaCmd.AddCommand(metaStringsCmd)
}

func printMeta(m meta.Result) {
	for _, key := range meta.Keys {
		if val, ok := m[key]; ok {
			fmt.Printf("%s: %s\n", key, val)
		} else {
			fmt.Printf("%s: (not found)\n", key)
		}
	}
}
