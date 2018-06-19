package net

import (
	"log"

	"github.com/Shopify/sarama"

	cluster "github.com/bsm/sarama-cluster"
)

// KafkaConsumer structure
type KafkaConsumer struct {
	config   *cluster.Config
	brokers  []string
	topics   []string
	consumer *cluster.Consumer
}

// NewKafkaConsumer return a new Kafka consumer
func NewKafkaConsumer(brokers []string, topics []string) *KafkaConsumer {
	config := cluster.NewConfig()
	config.Consumer.Return.Errors = false
	config.Group.Return.Notifications = false
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	// config.Net.KeepAlive = time.Duration(30) * time.Second
	config.ChannelBufferSize = 10000
	return &KafkaConsumer{brokers: brokers, topics: topics, config: config}
}

// Start the consumer
func (conn *KafkaConsumer) Start(name string) error {
	consumer, err := cluster.NewConsumer(conn.brokers, name, conn.topics, conn.config)
	conn.consumer = consumer
	if err == nil {
		go handleErrors(consumer)
		go handleNotifications(consumer)
	}
	return err
}

// GetMessageChan returns the message channel
func (conn *KafkaConsumer) GetMessageChan() <-chan *sarama.ConsumerMessage {
	return conn.consumer.Messages()
}

// MarkOffset for the given message
func (conn *KafkaConsumer) MarkOffset(msg *sarama.ConsumerMessage, meta string) {
	conn.consumer.MarkOffset(msg, meta)
}

// Close the consumer connection
func (conn *KafkaConsumer) Close() error {
	return conn.consumer.Close()
}

func handleErrors(consumer *cluster.Consumer) {
	for err := range consumer.Errors() {
		log.Printf("Error: %s\n", err.Error())
	}
}

func handleNotifications(consumer *cluster.Consumer) {
	for ntf := range consumer.Notifications() {
		log.Printf("Rebalanced: %+v\n", ntf)
	}
}
