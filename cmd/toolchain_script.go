package cmd

import (
	"fmt"
	"strings"

	"github.com/aura-studio/dynamic-cli/toolchain"
	"github.com/spf13/cobra"
)

var toolchainScriptCmd = &cobra.Command{
	Use:     "script",
	Short:   "导出 bash 环境变量脚本",
	Long:    "输出一个 bash 脚本到 stdout，用于设置环境变量 DYNAMIC_OS、DYNAMIC_ARCH、DYNAMIC_COMPILER。",
	Example: "  dynamic toolchain script | bash\n  eval \"$(dynamic toolchain script)\"\n",
	Args:    cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		v := toolchain.Detect()
		fmt.Print("#!/usr/bin/env bash\n")
		fmt.Print("set -euo pipefail\n")
		fmt.Printf("export DYNAMIC_OS=%s\n", bashSingleQuote(v.OS))
		fmt.Printf("export DYNAMIC_ARCH=%s\n", bashSingleQuote(v.Arch))
		fmt.Printf("export DYNAMIC_COMPILER=%s\n", bashSingleQuote(v.Compiler))
	},
}

func init() {
	toolchainCmd.AddCommand(toolchainScriptCmd)
}

func bashSingleQuote(s string) string {
	// Safe bash single-quote escaping: ' => '\''
	return "'" + strings.ReplaceAll(s, "'", "'\\''") + "'"
}
