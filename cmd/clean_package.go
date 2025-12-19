package cmd

import (
	"fmt"
	"os"

	"github.com/aura-studio/dynamic-cli/clean"
	"github.com/aura-studio/dynamic-cli/config"
	"github.com/spf13/cobra"
)

var cleanPackageCmd = &cobra.Command{
	Use:     "package",
	Short:   "Remove one package output directory",
	Long:    "Reads dynamic.yaml, locates the output directory for the given procedure, and removes that directory.",
	Example: "  dynamic clean package -c ./dynamic.yaml -p brazil\n",
	Args:    cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		cfgPath := resolveConfigPath(cmd)

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
		clean.CleanForProcedure(proc, clean.CleanTypePackage)
	},
}

func init() {
	cleanCmd.AddCommand(cleanPackageCmd)
	cleanPackageCmd.Flags().StringP("config", "c", "", "path to dynamic.yaml (default: ./dynamic.yaml if exists)")
	cleanPackageCmd.Flags().StringP("procedure", "p", "", "procedure name to clean (required)")
}
