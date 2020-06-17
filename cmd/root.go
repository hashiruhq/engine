package cmd

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// LogLevel Flag
var LogLevel = "info"
var LogFormat = "json"
var cfgFile string
var rootCmd = &cobra.Command{
	Use:   "trading_engine",
	Short: "Trading engine (name WIP) is a finantial market engine",
	Long: `A fast and flexible trading engine for the finantial market 
created by Around25 to support high frequency trading on crypto markets.
For a complete documentation and available licenses please contact https://around25.com`,
}

func init() {
	// set log level
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	initLoggingEnv()
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./.engine.yaml)")
	rootCmd.PersistentFlags().StringVarP(&LogLevel, "log-level", "", "info", "logging level to show (options: debug|info|warn|error|fatal|panic, default: info)")
	rootCmd.PersistentFlags().StringVarP(&LogFormat, "log-format", "", "info", "log format to generate (Options: json|pretty, default: json)")
	viper.SetConfigName(".engine")
	viper.AddConfigPath(".")                    // First try to load the config from the current directory
	viper.AddConfigPath("$HOME")                // Then try to load it from the HOME directory
	viper.AddConfigPath("/etc/trading_engine/") // As a last resort try to load it from the /etc/trading_engine
}

func initLoggingEnv() {
	// load log level from env by default
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel != "" {
		LogLevel = logLevel
	}
	// load log format from env by default
	logFormat := os.Getenv("LOG_FORMAT")
	if logFormat != "" {
		LogFormat = logFormat
	}
}

func initConfig() {
	// Don't forget to read config either from cfgFile, from current directory or from home directory!
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	}

	customizeLogger()
	// viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal().Err(err).Msg("Can't read configuration file")
	}
}

// Execute the commands
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal().Err(err).Msg("Unable to execute command")
	}
}

func customizeLogger() {
	if LogFormat == "pretty" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
	switch LogLevel {
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case "fatal":
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	case "panic":
		zerolog.SetGlobalLevel(zerolog.PanicLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
}
