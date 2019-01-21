package cmd

import (
	"log"

	"gitlab.com/around25/products/matching-engine/server"

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
		log.Printf("Loading server configuration...\n")
		if viper.ConfigFileUsed() != "" {
			log.Printf("Configuration file found at '%s'\n", viper.ConfigFileUsed())
		}
		cfg := server.LoadConfig(viper.GetViper())
		// start a new server
		log.Printf("Starting new server instance\n")
		srv := server.NewServer(cfg)
		// listen for new messages
		log.Printf("Listening for incomming orders...\n")
		srv.Listen()
	},
}
