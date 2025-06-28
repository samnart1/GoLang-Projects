package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use: "go-simple-logger",
	Short: "A simple and powerful logging utility",
	Long: `Go simple logger is a pro-grade logging utility that provides flexible message logging with timestamps, multiple output formats, and both CLI and HTTp server interfaces.
	
		Features:
			- Multiple log levels (DEBUG, INFO, WARN, ERROR, FATAL)
			- Flexible output (file, console, or both)
			- Customizable timestamp formats
			- Http server for remote logging
			- Real-time log monitoring
			- Log rotation support`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.go-simple-logger.yaml)")
	rootCmd.PersistentFlags().String("log-file", "", "log file path")
	rootCmd.PersistentFlags().String("log-level", "info", "log level (debug, info, warn, error, fatal)")
	rootCmd.PersistentFlags().String("log-format", "text", "log format (text, json)")
	rootCmd.PersistentFlags().String("timestamp-format", "2006-01-02 15:04:05", "timestamp format")

	viper.BindPFlag("log.file", rootCmd.PersistentFlags().Lookup("log-file"))
	viper.BindPFlag("log.level", rootCmd.PersistentFlags().Lookup("log-level"))
	viper.BindPFlag("log.format", rootCmd.PersistentFlags().Lookup("log-format"))
	viper.BindPFlag("log.timestamp_format", rootCmd.PersistentFlags().Lookup("timestamp-format"))
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home)
		viper.AddConfigPath(".")
		viper.SetConfigType("yaml")
		viper.SetConfigName(".go-simple-logger")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}