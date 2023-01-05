/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"
	"strings"

	"github.com/aura-studio/dynamic-cli/builder"
	"github.com/aura-studio/dynamic-cli/config"
	"github.com/spf13/cobra"
)

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build *.so in warehouse path",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if warehouse, err := cmd.Flags().GetString("warehouse"); err != nil {
			log.Panic(err)
		} else if warehouse != "" {
			config.SetDefaultWareHouse(warehouse)
		}

		if gover, err := cmd.Flags().GetString("gover"); err != nil {
			log.Panic(err)
		} else if gover != "" {
			config.SetDefaultGoVer(gover)
		}

		if debug, err := cmd.Flags().GetBool("debug"); err != nil {
			log.Panic(err)
		} else if debug {
			config.SetDefaultDebug(debug)
		}

		if len(args) > 0 {
			if strings.Contains(args[0], "@") {
				builder.BuildFromRepo(args[0], args[1:]...)
				return
			} else {
				builder.BuildFromJSONDir(args[0])
				return
			}
		}

		if file, err := cmd.Flags().GetString("file"); err != nil {
			log.Panic(err)
		} else if file != "" {
			builder.BuildFromJSONFile(file)
			return
		}

		if dir, err := cmd.Flags().GetString("dir"); err != nil {
			log.Panic(err)
		} else if dir != "" {
			builder.BuildFromJSONDir(dir)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)
	buildCmd.Flags().StringP("file", "f", "/tmp/dynamic.json", "path of config file")
	buildCmd.Flags().StringP("dir", "d", "/tmp", "path of config dir")
	buildCmd.Flags().StringP("warehouse", "w", "/tmp/warehouse", "path of warehouse")
	buildCmd.Flags().StringP("gover", "v", "1.18", "version of golang")
	buildCmd.Flags().BoolP("debug", "g", false, "build debug version")
}
