/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/aura-studio/dynamic-cli/config"
	"github.com/aura-studio/dynamic-cli/push"
	"github.com/spf13/cobra"
)

// pushCmd represents the push command
var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "Push using dynamic.yaml and specified procedure",
	Long:  `Reads dynamic.yaml and the given --procedure, then prepares push tasks (printing summary for now).`,
	Run: func(cmd *cobra.Command, args []string) {
		// resolve dynamic.yaml path: --config > current directory
		cfgPath, err := cmd.Flags().GetString("config")
		if err != nil {
			fmt.Println("error:", err)
			os.Exit(1)
		}
		if cfgPath == "" {
			cfgPath = filepath.Join(".", "dynamic.yaml")
		}

		// required procedure name
		proc, err := cmd.Flags().GetString("procedure")
		if err != nil {
			fmt.Println("error:", err)
			os.Exit(1)
		}
		if proc == "" {
			fmt.Println("error: procedure is required")
			os.Exit(1)
		}

		// parse and validate
		c := config.Parse(cfgPath)
		config.Validate(c)

		// build object based on procedure
		// compose procedure and call push entry
		procObj := config.CreateProcedure(c, proc)
		push.PushForProcedure(procObj)
	},
}

func init() {
	rootCmd.AddCommand(pushCmd)
	pushCmd.Flags().StringP("config", "c", "", "path to dynamic.yaml (default: ./dynamic.yaml)")
	pushCmd.Flags().StringP("procedure", "p", "", "procedure name to push (required)")
}
