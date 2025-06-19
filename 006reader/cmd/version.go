package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/samnart1/GoLang/006reader/pkg/version"
)

var (
	buildVersion = "dev"
	buildCommit	 = "none"
	buildDate	 = "unknown"
)

var versionCmd = &cobra.Command{
	Use: "version",
	Short: "Show version information",
	Long: "Display version, commit hash, and build date information",
	Run: runVersion,
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

func runVersion(cmd *cobra.Command, args []string) {
	info := version.GetInfo()

	fmt.Printf("Go File Reader\n")
	fmt.Printf("Version:	%s\n", info.Version)
	fmt.Printf("Commit:		%s\n", info.Commit)
	fmt.Printf("Build Date:	%s\n", info.Date)
	fmt.Printf("Go Version:	%s\n", info.GoVersion)
	fmt.Printf("Platform:	%s\n", info.OS, info.Arch)
}

func SetVersion(version, commit, date string) {
	buildVersion = version
	buildCommit = commit
	buildDate = date

	version.SetBuildInfo(version, commit, date)
}