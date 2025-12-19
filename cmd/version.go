package cmd

import (
	"fmt"
	"runtime"
	"runtime/debug"

	"github.com/spf13/cobra"
)

// These can be overridden via -ldflags "-X github.com/aura-studio/dynamic-cli/cmd.Version=v1.3.10"
var (
	Version = "dev"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		// If Version not injected via ldflags, try module build info (works with `go install module@version`).
		if Version == "dev" {
			if bi, ok := debug.ReadBuildInfo(); ok && bi != nil {
				if bi.Main.Version != "" && bi.Main.Version != "(devel)" {
					Version = bi.Main.Version
				}
			}
		}
		fmt.Printf("Version: %s\n", Version)
		fmt.Printf("Go:     %s %s/%s\n", runtime.Version(), runtime.GOOS, runtime.GOARCH)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
