/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/aura-studio/dynamic-cli/clean"
	"github.com/aura-studio/dynamic-cli/config"
	"github.com/spf13/cobra"
)

// cleanCmd represents the clean command
var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean using dynamic.yaml and specified procedure",
	Long:  `Reads dynamic.yaml and the given --procedure, then resolves paths to clean (printing summary for now).`,
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

		// compose procedure and call clean entry
		procObj := config.CreateProcedure(c, proc)
		clean.CleanForProcedure(procObj)
		fmt.Println("clean: invoked CleanForProcedure")
	},
}

func init() {
	rootCmd.AddCommand(cleanCmd)
	cleanCmd.Flags().String("config", "", "path to dynamic.yaml (default: ./dynamic.yaml)")
	cleanCmd.Flags().String("procedure", "", "procedure name to clean (required)")
}
