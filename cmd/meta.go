package cmd

import (
	"github.com/spf13/cobra"
)

var metaCmd = &cobra.Command{
	Use:   "meta",
	Short: "Read meta information from a shared library",
	Long: `Extract build metadata embedded in a .so file.

The meta command provides multiple methods to read metadata (module, version, 
built, os, arch, compiler, variant) that was injected during build time.

Available subcommands:
  strings   - Use 'strings' command to search for metadata patterns (fast, heuristic)
  nm        - Use 'go tool nm' to locate symbol addresses (shows symbol info)
  objdump   - Use 'objdump' to inspect .rodata section (low-level, verbose)
  call      - Load the .so and call the exported meta function (most reliable)

Example:
  dynamic meta strings libcgo_xxx.so
  dynamic meta nm libcgo_xxx.so
  dynamic meta objdump libcgo_xxx.so
  dynamic meta call libcgo_xxx.so`,
}

func init() {
	rootCmd.AddCommand(metaCmd)
}
