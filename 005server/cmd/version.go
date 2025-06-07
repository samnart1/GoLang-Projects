package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use: "version",
	Short: "print the version information",
	Long: `print the version information for the http server`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("go http server v1.0.0\n")
		fmt.Printf("build with %s\n", runtime.Version())
		fmt.Printf("OS/Arch: %s/%s\n", runtime.GOOS, runtime.GOARCH)
	},
}