package cmd

import (
	"fmt"
	"os"

	"github.com/aura-studio/dynamic-cli/clean"
	"github.com/aura-studio/dynamic-cli/config"
	"github.com/spf13/cobra"
)

var cleanUselessCmd = &cobra.Command{
	Use:     "useless",
	Short:   "Remove useless artifacts",
	Long:    "Reads dynamic-cli.yaml and deletes all files except .so and .json without date suffix under warehouse.local.",
	Example: "  dynamic clean useless -c ./dynamic-cli.yaml\n  dynamic clean useless -c ./dynamic-cli.yaml -p brazil\n",
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
			procName = c.Procedures[0].Name
		}
		proc := config.CreateProcedure(c, procName)
		clean.CleanForProcedure(proc, clean.CleanTypeUseless)
	},
}

func init() {
	cleanCmd.AddCommand(cleanUselessCmd)
}
