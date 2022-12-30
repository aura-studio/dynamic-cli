/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"
	"strings"

	"github.com/aura-studio/dynamic-cli/builder"
	"github.com/spf13/cobra"
)

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if warehouse, err := cmd.Flags().GetString("warehouse"); err != nil {
			log.Panic(err)
		} else if warehouse != "" {
			builder.DefaultConfig.WareHouse = warehouse
		}

		if gover, err := cmd.Flags().GetString("gover"); err != nil {
			log.Panic(err)
		} else if gover != "" {
			builder.DefaultConfig.GoVer = gover
		}

		if debug, err := cmd.Flags().GetBool("debug"); err != nil {
			log.Panic(err)
		} else if debug {
			builder.DefaultConfig.Debug = debug
		}

		if len(args) > 0 {
			if strings.Contains(args[0], "@") {
				builder.BuildFromRemote(args[0], args[1:]...)
				return
			} else {
				builder.BuildFromJSONPath(args[0])
			}
		}
		if config, err := cmd.Flags().GetString("config"); err != nil {
			log.Panic(err)
		} else if config != "" {
			builder.BuildFromJSONFile(cmd.Flag("config").Value.String())
		}

		if path, err := cmd.Flags().GetString("path"); err != nil {
			log.Panic(err)
		} else if path != "" {
			builder.BuildFromJSONPath(cmd.Flag("path").Value.String())
		}
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// buildCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// buildCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	buildCmd.Flags().StringP("config", "c", "", "path of config file")
	buildCmd.Flags().StringP("path", "p", "", "path of directory of dynamic.json")
	buildCmd.Flags().StringP("warehouse", "w", "", "path of warehouse")
	buildCmd.Flags().StringP("gover", "g", "", "version of golang")
	buildCmd.Flags().BoolP("debug", "d", false, "build debug version")
}
