package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/aura-studio/dynamic-cli/check"
	"github.com/aura-studio/dynamic-cli/config"
	"github.com/spf13/cobra"
)

// checkCmd represents the check command
var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Check current environment against toolchain",
	Long:  "Reads dynamic.yaml and the given --procedure, then verifies current OS/Arch/toolchain support.",
	Run: func(cmd *cobra.Command, args []string) {
		// resolve dynamic.yaml path: --config > current directory
		cfgPath, err := cmd.Flags().GetString("config")
		if err != nil {
			fmt.Println("error:", err)
			os.Exit(1)
		}
		if cfgPath == "" {
			cfgPath = filepath.Join(".", "dynamic.yaml")
		}

		// required procedure name
		proc, err := cmd.Flags().GetString("procedure")
		if err != nil {
			fmt.Println("error:", err)
			os.Exit(1)
		}
		if proc == "" {
			fmt.Println("error: procedure is required")
			os.Exit(1)
		}

		// parse and validate
		c := config.Parse(cfgPath)
		config.Validate(c)

		// compose procedure and run check
		procObj := config.CreateProcedure(c, proc)
		check.CheckForProcedure(procObj)
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)
	checkCmd.Flags().StringP("config", "c", "", "path to dynamic.yaml (default: ./dynamic.yaml)")
	checkCmd.Flags().StringP("procedure", "p", "", "procedure name to check (required)")
}
