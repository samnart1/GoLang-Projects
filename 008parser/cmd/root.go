package cmd

import (
	"fmt"
	"os"

	"github.com/samnart1/golang/008parser/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	cfg		*config.Config
)

var rootCmd = &cobra.Command{
	Use: "json-parser",
	Short: "A professional JSON parser with validation and formatting",
	Long: `A comprehensive JSON parser that can parse, validate and format JSON files
	
		Features:
		- Parse JSON files with multiple output formats
		- Validate JSON syntax and structure
		- Format/pretty-print JSON with customizable options
		- HTTP server for web-based parsing
		- Colored output for better readability`,

	Version: "1.0.0",
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.json-parser.yaml)")
	rootCmd.PersistentFlags().Bool("no-color", false, "disable colored output")
	rootCmd.PersistentFlags().Bool("verbose", false, "verbose output")
	rootCmd.PersistentFlags().Bool("quiet", false, "quiet mode (only errors)")

	viper.BindPFlag("no-color", rootCmd.PersistentFlags().Lookup("no-color"))
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
	viper.BindPFlag("quiet", rootCmd.PersistentFlags().Lookup("quiet"))
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
		viper.SetConfigName(".json-parser")
	}

	viper.AutomaticEnv()

	viper.SetDefault("DEFAULT_FORMAT", "pretty")
	viper.SetDefault("DEFAULT_INDENT", 2)
	viper.SetDefault("ENABLE_COLORS", true)
	viper.SetDefault("SERVER_HOST", "localhost")
	viper.SetDefault("SERVER_PORT", 8080)
	viper.SetDefault("MAX_FILE_SIZE", 10485760)
	viper.SetDefault("HTTP_TIMEOUT", 30)

	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}

	var err error
	cfg, err = config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}
}

func GetConfig() *config.Config {
	return cfg
}