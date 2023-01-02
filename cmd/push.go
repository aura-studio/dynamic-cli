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
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// pushCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	pushCmd.Flags().StringP("file", "f", "", "path of config file")
	pushCmd.Flags().StringP("dir", "d", "", "path of config dir")
	pushCmd.Flags().StringSliceP("remote", "r", nil, "remote warehouse")
}
