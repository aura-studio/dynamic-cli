/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"
	"strings"

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
		remotes, err := cmd.Flags().GetStringSlice("remote")
		if err != nil {
			log.Panic(err)
		}

		if len(args) > 0 {
			if strings.Contains(args[0], "@") {
				pusher.PushFromRepo(remotes, args[0], args[1:]...)
				return
			} else {
				pusher.PushFromJSONPath(remotes, args[0])
				return
			}
		}

		if file, err := cmd.Flags().GetString("file"); err != nil {
			log.Panic(err)
		} else if file != "" {
			pusher.PushFromJSONFile(remotes, file)
			return
		}

		if dir, err := cmd.Flags().GetString("dir"); err != nil {
			log.Panic(err)
		} else if dir != "" {
			pusher.PushFromJSONPath(remotes, dir)
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
	buildCmd.Flags().StringP("file", "f", "", "path of config file")
	buildCmd.Flags().StringP("dir", "d", "", "path of config dir")
	pushCmd.Flags().StringSliceP("remote", "r", nil, "remote warehouse")
}
