/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"
	"strings"

	"github.com/aura-studio/dynamic-cli/config"
	"github.com/aura-studio/dynamic-cli/pusher"
	"github.com/spf13/cobra"
)

// pushCmd represents the push command
var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "Push *.so to remote path",
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

		if remotes, err := cmd.Flags().GetStringSlice("remote"); err != nil {
			log.Panic(err)
		} else if len(remotes) > 0 {
			config.SetDefaultRemotes(remotes)
		}

		if len(args) > 0 {
			if strings.Contains(args[0], "@") {
				pusher.PushFromRepo(args[0], args[1:]...)
				return
			} else {
				pusher.PushFromJSONDir(args[0])
				return
			}
		}

		if file, err := cmd.Flags().GetString("file"); err != nil {
			log.Panic(err)
		} else if file != "" {
			pusher.PushFromJSONFile(file)
			return
		}

		if dir, err := cmd.Flags().GetString("dir"); err != nil {
			log.Panic(err)
		} else if dir != "" {
			pusher.PushFromJSONDir(dir)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(pushCmd)
	pushCmd.Flags().StringP("file", "f", "/tmp/dynamic.json", "path of config file")
	pushCmd.Flags().StringP("dir", "d", "/tmp", "path of config dir")
	pushCmd.Flags().StringP("warehouse", "w", "/tmp/warehouse", "path of warehouse")
	pushCmd.Flags().StringSliceP("remote", "r", nil, "remote warehouse")
}
