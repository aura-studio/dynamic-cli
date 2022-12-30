/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"

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
			log.Fatal(err)
		} else if len(remotes) > 0 {
			// fmt.Println(remotes)
		}

		if file, err := cmd.Flags().GetString("file"); err != nil {
			log.Panic(err)
		} else if file != "" {
			// pusher.BuildFromJSONFile(cmd.Flag("config").Value.String())
		}

		if path, err := cmd.Flags().GetString("path"); err != nil {
			log.Panic(err)
		} else if path != "" {
			// builder.BuildFromJSONPath(cmd.Flag("path").Value.String())
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
	pushCmd.Flags().StringSliceP("remote", "r", nil, "s3 remote URL")
	buildCmd.Flags().StringP("file", "f", "", "path of config file")
	buildCmd.Flags().StringP("dir", "d", "", "path of config dir")
}
