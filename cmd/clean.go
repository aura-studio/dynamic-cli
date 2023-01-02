/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"
	"strings"

	"github.com/aura-studio/dynamic-cli/cleaner"
	"github.com/aura-studio/dynamic-cli/config"
	"github.com/spf13/cobra"
)

// cleanCmd represents the clean command
var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean build files",
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

		all, err := cmd.Flags().GetBool("all")
		if err != nil {
			log.Panic(err)
		}

		if len(args) > 0 {
			if strings.Contains(args[0], "@") {
				cleaner.CleanFromRepo(all, args[0], args[1:]...)
				return
			} else {
				cleaner.CleanFromJSONDir(all, args[0])
				return
			}
		}

		if file, err := cmd.Flags().GetString("file"); err != nil {
			log.Panic(err)
		} else if file != "" {
			cleaner.CleanFromJSONFile(all, file)
			return
		}

		if dir, err := cmd.Flags().GetString("dir"); err != nil {
			log.Panic(err)
		} else if dir != "" {
			cleaner.CleanFromJSONDir(all, dir)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(cleanCmd)
	cleanCmd.Flags().StringP("file", "f", "", "path of config file")
	cleanCmd.Flags().StringP("dir", "d", "", "path of config dir")
	cleanCmd.Flags().StringP("warehouse", "w", "", "path of warehouse")
	cleanCmd.Flags().BoolP("all", "a", false, "clean all packages")
}
