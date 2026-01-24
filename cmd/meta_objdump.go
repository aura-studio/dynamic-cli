package cmd

import (
	"fmt"
	"os"

	"github.com/aura-studio/dynamic-cli/meta"
	"github.com/spf13/cobra"
)

var metaObjdumpCmd = &cobra.Command{
	Use:   "objdump <file.so>",
	Short: "Extract meta using objdump",
	Long: `Use 'objdump' to inspect the .rodata section of the .so file.

This method dumps the read-only data section where string constants are stored.
The output is raw hex + ASCII, useful for low-level debugging but verbose.
You'll need to manually search for metadata strings in the output.

The command executed is:
  objdump -s -j .rodata <file.so>

For a more readable output, consider piping through grep:
  dynamic meta objdump libcgo_xxx.so | grep -A2 "github.com"

Example:
  dynamic meta objdump libcgo_xxx.so`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		soPath := args[0]
		if _, err := os.Stat(soPath); os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "file not found: %s\n", soPath)
			os.Exit(1)
		}

		out, err := meta.ReadObjdump(soPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to read meta: %v\n", err)
			os.Exit(1)
		}

		fmt.Print(out)
	},
}

func init() {
	metaCmd.AddCommand(metaObjdumpCmd)
}
