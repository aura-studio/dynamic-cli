package cmd

import "github.com/spf13/cobra"

// toolchainCmd groups toolchain-related subcommands.
var toolchainCmd = &cobra.Command{
	Use:   "toolchain",
	Short: "toolchain 检测与查询",
	Long:  "提供 toolchain 的检测(check)与查询(describe)能力。",
}

func init() {
	rootCmd.AddCommand(toolchainCmd)
}
