package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/hashiruhq/engine/model"

	"github.com/spf13/cobra"
)

func init() {
	dumpOrderbook.Flags().StringVarP(&backupPath, "path", "p", "./backups", "Path to the backup files")
	dumpOrderbook.Flags().StringVarP(&market, "market", "m", "", "The name of the market to dump")
	dumpOrderbook.Flags().StringVarP(&format, "format", "f", "json", "Output format")
	rootCmd.AddCommand(dumpOrderbook)
}

var backupPath string = ""
var market string = ""
var format string = "json"

var dumpOrderbook = &cobra.Command{
	Use:   "dump_orderbook",
	Short: "Dump a list of all orders from a backup file",
	Long:  `Dump a list of all orders from a backup file`,
	Run: func(cmd *cobra.Command, args []string) {
		if market == "" {
			log.Fatalln("Please provide the market id to dump using the --market or -m flag")
		}
		dump_orderbook(backupPath, market, format)
	},
}

func dump_orderbook(path, market, format string) {
	file := filepath.Join(path, market)
	file = filepath.Clean(file)
	content, err := os.ReadFile(file)
	if err != nil {
		log.Fatalln("Unable to import market backup from: ", file)
		return
	}

	var backup model.MarketBackup
	backup.FromBinary(content)
	output, err := json.MarshalIndent(backup, "", "  ")
	if err != nil {
		log.Println("Unable to generate JSON file for backup: ", err)
		return
	}
	fmt.Println(string(output))
}
