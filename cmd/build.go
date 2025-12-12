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

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build using dynamic.yaml and specified procedure",
	Long:  `Reads dynamic.yaml and the given --procedure, then constructs a Build object.`,
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

		// build object
		b := config.BuildForProcedure(c, proc)
		// For now, just print summary; integration with builder can follow
		fmt.Printf("Build plan:\nToolchain: %s/%s %s (%s)\nWarehouse: %s -> %v\nSource: %s %s@%s\nTarget: %s %s@%s\n",
			b.Toolchain.OS, b.Toolchain.Arch, b.Toolchain.Compiler, b.Toolchain.Variant,
			b.Warehouse.Local, b.Warehouse.Remote,
			b.Source.Repo, b.Source.Module, b.Source.Version,
			b.Target.Namespace, b.Target.Package, b.Target.Version,
		)
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)
	buildCmd.Flags().StringP("config", "c", "", "path to dynamic.yaml (default: ./dynamic.yaml)")
	buildCmd.Flags().StringP("procedure", "p", "", "procedure name to build (required)")
}
