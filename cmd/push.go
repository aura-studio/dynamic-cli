/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
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
			log.Panic(err)
		}
		if cfgPath == "" {
			cfgPath = filepath.Join(".", "dynamic.yaml")
		}

		// required procedure name
		proc, err := cmd.Flags().GetString("procedure")
		if err != nil {
			log.Panic(err)
		}
		if proc == "" {
			log.Panic("procedure is required")
		}

		// parse and validate
		c := config.Parse(cfgPath)
		config.Validate(c)

		// build object based on procedure
		// compose procedure and call push entry
		procObj := config.CreateProcedure(c, proc)
		push.PushForProcedure(procObj)
		fmt.Println("push: invoked PushForProcedure")
	},
}

func init() {
	rootCmd.AddCommand(pushCmd)
	pushCmd.Flags().String("config", "", "path to dynamic.yaml (default: ./dynamic.yaml)")
	pushCmd.Flags().String("procedure", "", "procedure name to push (required)")
}
