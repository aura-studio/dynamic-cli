/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/aura-studio/dynamic-cli/clean"
	"github.com/aura-studio/dynamic-cli/config"
	"github.com/spf13/cobra"
)

// cleanCmd represents the clean command
var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean using dynamic.yaml and specified procedure",
	Long:  `Reads dynamic.yaml and the given --procedure, then resolves paths to clean (printing summary for now).`,
	Run: func(cmd *cobra.Command, args []string) {
		// resolve dynamic.yaml path: --config > current directory
		cfgPath, err := cmd.Flags().GetString("config")
		if err != nil {
			log.Panic(err)
		}
		if cfgPath == "" {
			cfgPath = filepath.Join(".", "dynamic.yaml")
		}

		// clean type
		t, err := cmd.Flags().GetString("type")
		if err != nil {
			log.Panic(err)
		}
		if t == "" {
			log.Panic("clean type is required: cache|package|all")
		}
		var ct clean.CleanType
		switch t {
		case "cache":
			ct = clean.CleanTypeCache
		case "package":
			ct = clean.CleanTypePackage
		case "all":
			ct = clean.CleanTypeAll
		default:
			log.Panic("invalid clean type: " + t)
		}

		// procedure required when type=package
		proc, err := cmd.Flags().GetString("procedure")
		if err != nil {
			log.Panic(err)
		}
		if ct == clean.CleanTypePackage && proc == "" {
			log.Panic("procedure is required when type=package")
		}

		// parse and validate
		c := config.Parse(cfgPath)
		config.Validate(c)

		// compose procedure when needed and call clean entry
		var procObj config.Procedure
		if ct == clean.CleanTypePackage || ct == clean.CleanTypeCache {
			procObj = config.CreateProcedure(c, proc)
		} else {
			// for all, a dummy procedure with warehouse local is sufficient
			// but CreateProcedure enforces name; use provided or pick first
			if proc == "" && len(c.Procedures) > 0 {
				procObj = config.CreateProcedure(c, c.Procedures[0].Name)
			} else if proc != "" {
				procObj = config.CreateProcedure(c, proc)
			}
		}
		clean.CleanForProcedure(procObj, ct)
		fmt.Println("clean: done")
	},
}

func init() {
	rootCmd.AddCommand(cleanCmd)
	cleanCmd.Flags().String("config", "", "path to dynamic.yaml (default: ./dynamic.yaml)")
	cleanCmd.Flags().String("procedure", "", "procedure name (required when type=package)")
	cleanCmd.Flags().String("type", "", "clean type: cache|package|all (required)")
}
