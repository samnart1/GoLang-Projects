package cmd

import (
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

func initConfig() {}