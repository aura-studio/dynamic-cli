package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/aura-studio/dynamic-cli/config"
	"github.com/aura-studio/dynamic-cli/toolchain"
	"github.com/spf13/cobra"
)

var toolchainCheckCmd = &cobra.Command{
	Use:     "check",
	Short:   "一次性校验 OS / Arch / Compiler",
	Long:    "读取 dynamic.yaml 中指定 procedure 的 environment.toolchain 配置，并一次性校验 OS/Arch/Compiler 三个字段。任一字段不匹配则退出码为 1。",
	Example: "  dynamic toolchain check -c ./dynamic.yaml -p brazil\n",
	Args:    cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		cfgPath, err := cmd.Flags().GetString("config")
		if err != nil {
			fmt.Println("error:", err)
			os.Exit(1)
		}
		if cfgPath == "" {
			cfgPath = filepath.Join(".", "dynamic.yaml")
		}

		procName, err := cmd.Flags().GetString("procedure")
		if err != nil {
			fmt.Println("error:", err)
			os.Exit(1)
		}
		if procName == "" {
			fmt.Println("error: procedure is required")
			os.Exit(1)
		}

		c := config.Parse(cfgPath)
		config.Validate(c)
		proc := config.CreateProcedure(c, procName)

		expected := toolchain.Values{
			OS:       proc.Toolchain.OS,
			Arch:     proc.Toolchain.Arch,
			Compiler: proc.Toolchain.Compiler,
		}
		if ok := toolchain.Check(expected); !ok {
			os.Exit(1)
		}
	},
}

func init() {
	toolchainCmd.AddCommand(toolchainCheckCmd)
	toolchainCheckCmd.Flags().StringP("config", "c", "", "path to dynamic.yaml (default: ./dynamic.yaml)")
	toolchainCheckCmd.Flags().StringP("procedure", "p", "", "procedure name to check (required)")
}
