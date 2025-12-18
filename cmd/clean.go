/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import "github.com/spf13/cobra"

// cleanCmd represents the clean command
var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "清理构建产物",
	Long:  "清理 warehouse 下的构建产物，支持 cache/package/all 三种子命令。支持 -c 指定配置；不传 -c 时仅在当前目录存在 dynamic.yaml 才会自动使用。",
}

func init() {
	rootCmd.AddCommand(cleanCmd)
}
