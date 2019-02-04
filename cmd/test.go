package cmd

import (
	"log"
	"math"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gitlab.com/around25/products/matching-engine/engine"
	"gitlab.com/around25/products/matching-engine/net"
	"gitlab.com/around25/products/matching-engine/server"

	"github.com/Shopify/sarama"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	testCmd.Flags().StringVarP(&testCase, "case", "c", "gen_orders", "Generate orders in the configured consumers.")
	testCmd.Flags().IntVarP(&timeout, "timeout", "t", 10, "The timeout until the generator should exit")
	testCmd.Flags().IntVarP(&delay, "delay", "d", 0, "The delay between orders")
	testCmd.Flags().IntVarP(&topicCount, "topic_count", "", 0, "The maximum number of topics to generate orders for. Default: all topics configured.")
	rootCmd.AddCommand(testCmd)
}

var testCase string
var timeout int
var delay int
var topicCount int
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Test the trading engine by generating fake orders",
	Long:  `Test the system by generating test orders and returning the number of trades completed`,
	Run: func(cmd *cobra.Command, args []string) {
		switch testCase {
		case "gen_orders":
			test_gen_orders(timeout, delay, topicCount)
		}
	},
}

func test_gen_orders(timeout, delay, topicCount int) {
	rand.Seed(42)
	// load server configuration from server
	cfg := server.LoadConfig(viper.GetViper())

	producers := make([]net.KafkaProducer, 0, len(cfg.Brokers.Consumers))
	for _, consumer := range cfg.Brokers.Consumers {
		producer := net.NewKafkaAsyncProducer(consumer.Hosts)
		producer.Start()
		go func(producer net.KafkaProducer, topics []string) {
			maxTopics := len(topics)
			if topicCount > 0 {
				maxTopics = topicCount
			}

			index := 0
			for {
				topicIndex := rand.Intn(maxTopics)
				id := index
				price := uint64(math.Floor((1 + float64(index) + 10000*rand.Float64()) * 100000000))
				amount := uint64(math.Floor((10001 - 10000*rand.Float64()) * 100000000))
				side := engine.MarketSide_Buy
				if rand.Intn(2)%2 == 1 {
					side = engine.MarketSide_Sell
				}
				order := engine.Order{
					ID:        uint64(id),
					Price:     price,
					Amount:    amount,
					Side:      side,
					EventType: engine.CommandType_NewOrder,
					Type:      engine.OrderType_Limit,
				}
				data, _ := order.ToBinary()
				producer.Input() <- &sarama.ProducerMessage{
					Topic: topics[topicIndex],
					Value: sarama.ByteEncoder(data),
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
