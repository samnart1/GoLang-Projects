package cmd

import "github.com/spf13/cobra"

var configCmd = &cobra.Command{
	Use: "config",
	Short: "Configure game settings",
	Long: "View and modify game configuration settings",
	RunE: runConfig,
}

func runConfig(cmd *cobra.Command, args []string) error {
	// cfg, err := config.Load()
	// if err != nil {
	// 	return err
	// }

	// ui.ShowConfig(cfg)
	
	return nil
}