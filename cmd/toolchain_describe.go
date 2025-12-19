package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/aura-studio/dynamic-cli/toolchain"
	"github.com/spf13/cobra"
)

var toolchainDescribeCmd = &cobra.Command{
	Use:     "describe [os|arch|compiler|all]",
	Short:   "Print current toolchain values",
	Long:    "Detects and prints toolchain values from the local machine. For a single field it prints only the value; for 'all' it prints three labeled lines (OS/Arch/Compiler).",
	Example: "  dynamic toolchain describe os\n  dynamic toolchain describe arch\n  dynamic toolchain describe compiler\n  dynamic toolchain describe all\n",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		kind := strings.ToLower(strings.TrimSpace(args[0]))
		switch kind {
		case string(toolchain.KindOS):
			fmt.Println(toolchain.Describe(toolchain.KindOS))
		case string(toolchain.KindArch):
			fmt.Println(toolchain.Describe(toolchain.KindArch))
		case string(toolchain.KindCompiler):
			fmt.Println(toolchain.Describe(toolchain.KindCompiler))
		case "all":
			fmt.Printf("OS: %s\n", toolchain.Describe(toolchain.KindOS))
			fmt.Printf("Arch: %s\n", toolchain.Describe(toolchain.KindArch))
			fmt.Printf("Compiler: %s\n", toolchain.Describe(toolchain.KindCompiler))
		default:
			fmt.Println("error: kind must be one of: os, arch, compiler, all")
			os.Exit(1)
		}
	},
}

func init() {
	toolchainCmd.AddCommand(toolchainDescribeCmd)
}
