/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/aura-studio/dynamic-cli/config"
	"github.com/aura-studio/dynamic-cli/push"
	"github.com/spf13/cobra"
)

// pushCmd represents the push command
var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "Push using dynamic-cli.yaml and specified procedure",
	Long:  `Reads dynamic-cli.yaml and the given --procedure, then prepares push tasks (printing summary for now).`,
	Run: func(cmd *cobra.Command, args []string) {
		cfgPath := resolveConfigPath(cmd)

		// procedure name is optional
		proc, err := cmd.Flags().GetString("procedure")
		if err != nil {
			fmt.Println("error:", err)
			os.Exit(1)
		}

		// parse and validate
		c := config.Parse(cfgPath)
		config.Validate(c)

		if proc == "" {
			// push all procedures
			fmt.Println("No procedure specified, pushing all procedures...")
			procedures := config.GetAllProcedures(c)
			for _, procName := range procedures {
				fmt.Printf("\nPushing procedure: %s\n", procName)
				procObj := config.CreateProcedure(c, procName)
				push.PushForProcedure(procObj)
			}
			fmt.Println("\nAll procedures pushed successfully.")
		} else {
			// push specified procedure
			procObj := config.CreateProcedure(c, proc)
			push.PushForProcedure(procObj)
		}
	},
}

func init() {
	rootCmd.AddCommand(pushCmd)
	pushCmd.Flags().StringP("config", "c", "", "path to dynamic-cli.yaml (default: ./dynamic-cli.yaml or ./dynamic-cli.yml)")
	pushCmd.Flags().StringP("procedure", "p", "", "procedure name to push (optional, pushes all if not specified)")
}
