package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/aura-studio/dynamic-cli/clean"
	"github.com/aura-studio/dynamic-cli/config"
	"github.com/spf13/cobra"
)

var cleanAllCmd = &cobra.Command{
	Use:     "all",
	Short:   "清理整个 warehouse",
	Long:    "读取 dynamic.yaml 并删除 warehouse.local 下的所有内容。可以通过 -p 指定 procedure 来选择 warehouse；不指定则使用配置中的第一个 procedure。",
	Example: "  dynamic clean all -c ./dynamic.yaml\n  dynamic clean all -c ./dynamic.yaml -p brazil\n",
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

		c := config.Parse(cfgPath)
		config.Validate(c)

		if procName == "" {
			// Validate ensures procedures is non-empty.
			procName = c.Procedures[0].Name
		}
		proc := config.CreateProcedure(c, procName)
		clean.CleanForProcedure(proc, clean.CleanTypeAll)
	},
}

func init() {
	cleanCmd.AddCommand(cleanAllCmd)
	cleanAllCmd.Flags().StringP("config", "c", "", "path to dynamic.yaml (default: ./dynamic.yaml)")
	cleanAllCmd.Flags().StringP("procedure", "p", "", "procedure name to select warehouse (optional)")
}
