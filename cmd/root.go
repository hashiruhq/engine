package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var rootCmd = &cobra.Command{
	Use:   "trading_engine",
	Short: "Trading engine (name WIP) is a finantial market engine",
	Long: `A fast and flexible trading engine for the finantial market 
created by Around25 to support high frequency trading on crypto markets.
For a complete documentation and available licenses please contact https://around25.com`,
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./.engine.yaml)")
	viper.SetConfigName(".engine")
	viper.AddConfigPath(".")                    // First try to load the config from the current directory
	viper.AddConfigPath("$HOME")                // Then try to load it from the HOME directory
	viper.AddConfigPath("/etc/trading_engine/") // As a last resort try to load it from the /etc/trading_engine
}

func initConfig() {
	// Don't forget to read config either from cfgFile, from current directory or from home directory!
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	}

	// viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Can't read config:", err)
		os.Exit(1)
	}
}

// Execute the commands
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
