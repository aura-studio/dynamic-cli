/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/aura-studio/dynamic-cli/build"
	"github.com/aura-studio/dynamic-cli/config"
	"github.com/spf13/cobra"
)

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build using dynamic.yaml and specified procedure",
	Long:  `Reads dynamic.yaml and the given --procedure, then constructs a Build object.`,
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
			// build all procedures
			fmt.Println("No procedure specified, building all procedures...")
			procedures := config.GetAllProcedures(c)
			for _, procName := range procedures {
				fmt.Printf("\nBuilding procedure: %s\n", procName)
				procObj := config.CreateProcedure(c, procName)
				build.BuildForProcedure(procObj)
			}
			fmt.Println("\nAll procedures built successfully.")
		} else {
			// build specified procedure
			procObj := config.CreateProcedure(c, proc)
			build.BuildForProcedure(procObj)
		}
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)
	buildCmd.Flags().StringP("config", "c", "", "path to dynamic.yaml (default: ./dynamic.yaml or ./dynamic.yml)")
	buildCmd.Flags().StringP("procedure", "p", "", "procedure name to build (optional, builds all if not specified)")
}
