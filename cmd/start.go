package cmd

import (
	"gitlab.com/around25/products/matching-engine/server"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(serverCmd)
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the trading engine and listen for new orders on the configured markets",
	Long:  `Connect to the configured message queue and listen for new requests and generate events like trades and market state`,
	Run: func(cmd *cobra.Command, args []string) {
		// load server configuration from server
		log.Debug().Str("section", "server").Str("action", "init").Msg("Loading server configuration")
		if viper.ConfigFileUsed() != "" {
			log.Debug().Str("section", "server").Str("action", "init").Str("path", viper.ConfigFileUsed()).Msg("Configuration file loaded")
		}
		cfg := server.LoadConfig(viper.GetViper())
		// start a new server
		log.Debug().Str("section", "server").Str("action", "init").Msg("Starting new server instance")
		srv := server.NewServer(cfg)
		// listen for new messages
		log.Info().Str("section", "server").Str("action", "init").Msg("Listening for incoming orders")
		srv.Listen()
	},
}
