package cmd

import (
	"fmt"
	"os"

	"github.com/aura-studio/dynamic-cli/meta"
	"github.com/spf13/cobra"
)

var metaNmCmd = &cobra.Command{
	Use:   "nm <file.so>",
	Short: "Extract meta using go tool nm",
	Long: `Use 'go tool nm' to locate Meta symbol addresses in the .so file.

This method lists symbol table entries and filters for the injected Meta 
variables (MetaModule, MetaVersion, MetaBuilt, etc.). It shows the memory
address and type of each symbol, confirming they exist in the binary.

Note: This method shows symbol presence and addresses but cannot directly
extract the string values. Use 'strings' or 'call' method for actual values.

Symbol names searched:
  main.MetaModule
  main.MetaVersion
  main.MetaBuilt
  main.MetaOS
  main.MetaArch
  main.MetaCompiler
  main.MetaVariant

Example:
  dynamic meta nm libcgo_xxx.so`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		soPath := args[0]
		if _, err := os.Stat(soPath); os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "file not found: %s\n", soPath)
			os.Exit(1)
		}

		found, err := meta.ReadNm(soPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to read meta: %v\n", err)
			os.Exit(1)
		}

		if len(found) == 0 {
			fmt.Println("No Meta symbols found in the binary.")
			return
		}

		fmt.Println("Meta symbols found:")
		symbols := []string{
			"main.MetaModule",
			"main.MetaVersion",
			"main.MetaBuilt",
			"main.MetaOS",
			"main.MetaArch",
			"main.MetaCompiler",
			"main.MetaVariant",
		}
		for _, sym := range symbols {
			if info, ok := found[sym]; ok {
				fmt.Printf("  %s\n", info)
			}
		}
	},
}

func init() {
	metaCmd.AddCommand(metaNmCmd)
}
