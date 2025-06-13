package cmd

import (
	"fmt"
	"os"

	"github.com/samnart1/GoLang/006reader/internal/config"
	"github.com/samnart1/GoLang/006reader/internal/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	verbose	bool
	cfg		*config.Config
	log		*logger.Logger
)

var rootCmd = &cobra.Command{
	Use: "go-file-reader",
	Short: "A powerful file reading utility",
	Long: `Go File Reader is a command line tool for reading and displaying file contents with various formatting options.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return initializeConfig()
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.go-file-reader.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")

	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home)
		viper.AddConfigPath(".")
		viper.SetConfigFile("yaml")
		viper.SetConfigFile(".go-file-reader")
	}

	viper.SetEnvPrefix("FILE_READER")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}

func initializeConfig() error {
	var err error

	cfg, err = config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	log, err = logger.New()
	if err != nil {
		return fmt.Errorf("failed to initialize logger: %w", err)
	}

	return nil
}