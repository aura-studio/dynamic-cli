package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/aura-studio/dynamic-cli/clean"
	"github.com/aura-studio/dynamic-cli/config"
	"github.com/spf13/cobra"
)

var cleanCacheCmd = &cobra.Command{
	Use:     "cache",
	Short:   "清理缓存（保留主产物）",
	Long:    "读取 dynamic.yaml 并定位到指定 procedure 的产物目录，清理目录内除主产物文件外的缓存内容。",
	Example: "  dynamic clean cache -c ./dynamic.yaml -p brazil\n",
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
		clean.CleanForProcedure(proc, clean.CleanTypeCache)
	},
}

func init() {
	cleanCmd.AddCommand(cleanCacheCmd)
	cleanCacheCmd.Flags().StringP("config", "c", "", "path to dynamic.yaml (default: ./dynamic.yaml)")
	cleanCacheCmd.Flags().StringP("procedure", "p", "", "procedure name to clean (required)")
}
