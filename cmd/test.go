package cmd

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"
	"trading_engine/net"
	"trading_engine/server"

	"github.com/Shopify/sarama"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	testCmd.Flags().StringVarP(&testCase, "case", "c", "gen_orders", "Generate orders in the configured consumers.")
	testCmd.Flags().IntVarP(&timeout, "timeout", "t", 10, "The timeout until the generator should exit")
	testCmd.Flags().IntVarP(&delay, "delay", "d", 0, "The delay between orders")
	rootCmd.AddCommand(testCmd)
}

var testCase string
var timeout int
var delay int
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Test the trading engine by generating fake orders",
	Long:  `Test the system by generating test orders and returning the number of trades completed`,
	Run: func(cmd *cobra.Command, args []string) {
		switch testCase {
		case "gen_orders":
			test_gen_orders(timeout, delay)
		}
	},
}

func test_gen_orders(timeout, delay int) {
	rand.Seed(42)
	// load server configuration from server
	cfg := server.LoadConfig(viper.GetViper())

	producers := make([]net.KafkaProducer, 0, len(cfg.Brokers.Consumers))
	for _, consumer := range cfg.Brokers.Consumers {
		producer := net.NewKafkaAsyncProducer(consumer.Hosts)
		producer.Start()
		go func(producer net.KafkaProducer, topics []string) {
			maxTopics := len(topics)
			index := 0
			for {
				topicIndex := rand.Intn(maxTopics)
				id := "ID_" + fmt.Sprintf("%d", index)
				price := 1 + 3*index + int(math.Ceil(10000*rand.Float64()))
				amount := 10001 - int(math.Ceil(10000*rand.Float64()))
				side := int8(1 + rand.Intn(2)%2)
				order := fmt.Sprintf(`{"base": "sym", "market": "tst", "id":"%s", "price": %d, "amount": %d, "side": %d, "category": 1}`, id, price, amount, side)
				producer.Input() <- &sarama.ProducerMessage{
					Topic: topics[topicIndex],
					Value: sarama.ByteEncoder(([]byte)(order)),
				}
				index++
				if delay > 0 {
					log.Println(topics[topicIndex], order)
					time.Sleep(time.Duration(delay) * time.Second)
				}
			}
		}(producer, consumer.Topics)
		go func(producer net.KafkaProducer) {
			errors := producer.Errors()
			for err := range errors {
				value, _ := err.Msg.Value.Encode()
				log.Print("Error received from trades producer ", (string)(value), err)
			}
		}(producer)
		producers = append(producers, producer)
	}

	// wait for close signal from the user before exiting
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt, syscall.SIGTERM)

	wait := time.After(time.Duration(timeout) * time.Second)
	select {
	case sig := <-sigc:
		log.Printf("Caught signal %s: Shutting down in 3 seconds...", sig)
		log.Println("Closing all test producers...")
	case <-wait:
		log.Printf("Timeout expired after %d seconds...", timeout)
		log.Println("Closing all test producers...")
	}

	for _, producer := range producers {
		producer.Close()
	}

	os.Exit(0)
}
