package cmd

import (
	"fmt"
	"os"

	"github.com/aura-studio/dynamic-cli/clean"
	"github.com/aura-studio/dynamic-cli/config"
	"github.com/spf13/cobra"
)

var cleanAllCmd = &cobra.Command{
	Use:     "all",
	Short:   "Remove the entire warehouse",
	Long:    "Reads dynamic-cli.yaml and deletes all contents under warehouse.local. You may pass -p to choose which procedure's warehouse to use; if omitted, the first procedure in config is used.",
	Example: "  dynamic clean all -c ./dynamic-cli.yaml\n  dynamic clean all -c ./dynamic-cli.yaml -p brazil\n",
	Args:    cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		cfgPath := resolveConfigPath(cmd)

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
	cleanAllCmd.Flags().StringP("config", "c", "", "path to dynamic-cli.yaml (default: ./dynamic-cli.yaml or ./dynamic-cli.yml)")
	cleanAllCmd.Flags().StringP("procedure", "p", "", "procedure name to select warehouse (optional)")
}
