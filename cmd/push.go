/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/aura-studio/dynamic-cli/config"
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
		b := config.BuildForProcedure(c, proc)
		// For now, just print push plan
		fmt.Printf("Push plan:\nWarehouse local: %s\nWarehouse remote: %v\nArtifact: namespace=%s package=%s version=%s\n",
			b.Warehouse.Local,
			b.Warehouse.Remote,
			b.Target.Namespace, b.Target.Package, b.Target.Version,
		)
	},
}

func init() {
	rootCmd.AddCommand(pushCmd)
	pushCmd.Flags().String("config", "", "path to dynamic.yaml (default: ./dynamic.yaml)")
	pushCmd.Flags().String("procedure", "", "procedure name to push (required)")
}
