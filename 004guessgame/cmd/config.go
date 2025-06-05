package cmd

import (
	"github.com/samnart1/GoLang-Projects/004guessgame/internal/config"
	"github.com/samnart1/GoLang-Projects/004guessgame/internal/ui"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use: "config",
	Short: "Configure game settings",
	Long: "View and modify game configuration settings",
	RunE: runConfig,
}

func runConfig(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	ui.ShowConfig(cfg)
	
	return nil
}