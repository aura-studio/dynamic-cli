package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/aura-studio/dynamic-cli/toolchain"
	"github.com/spf13/cobra"
)

var toolchainDescribeCmd = &cobra.Command{
	Use:   "describe [os|arch|compiler]",
	Short: "打印当前 toolchain 值",
	Long:  "从本机环境探测并输出指定字段的值。输出仅包含值本身。",
	Example: "  dynamic toolchain describe os\n  dynamic toolchain describe arch\n  dynamic toolchain describe compiler\n",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		kind := strings.ToLower(strings.TrimSpace(args[0]))
		switch kind {
		case string(toolchain.KindOS):
			fmt.Println(toolchain.Describe(toolchain.KindOS))
		case string(toolchain.KindArch):
			fmt.Println(toolchain.Describe(toolchain.KindArch))
		case string(toolchain.KindCompiler):
			fmt.Println(toolchain.Describe(toolchain.KindCompiler))
		default:
			fmt.Println("error: kind must be one of: os, arch, compiler")
			os.Exit(1)
		}
	},
}

func init() {
	toolchainCmd.AddCommand(toolchainDescribeCmd)
}
