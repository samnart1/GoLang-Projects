package cmd

import (
	"fmt"
	"os"

	"github.com/samnart1/golang/010timer/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	verbose bool
	quiet bool
	cfg *config.Config
)

var rootCmd = &cobra.Command{
	Use: "timer",
	Short: "A powerful command-line timer application",
	Long: `Go Simple Timer is a feature-rich command-line timer application that supports
		
		- Coundown timers with flexible duration formats
		- Stopwatch functionality with lap times
		- Alarms for specific times
		- Mulitiple notification types (sound, desktop, terminal)
		- Customizable configuration
		
		Examples: 
			timer countdown -d 5m
			timer countdown -d 30s --sound
			timer stopwatch
			timer alarm -t 14:30`,

	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.timer.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false, "quiet mode (no sound alerts)")

	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
	viper.BindPFlag("quiet", rootCmd.PersistentFlags().Lookup("quiet"))
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			 fmt.Errorf("error finding home directory: %w", err)
		}

		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".timer")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil && verbose {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}

	var err error
	cfg, err = config.Load()
	if err != nil {
		 fmt.Errorf("error loading config: %w", err)
	}
}

func GetConfig() *config.Config {
	return cfg
}