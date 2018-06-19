package net

import (
	"log"

	"github.com/Shopify/sarama"
)

// KafkaProducer structure
type KafkaProducer struct {
	brokers  []string
	topic    string
	producer sarama.SyncProducer
	config   *sarama.Config
}

// NewKafkaProducer returns a new producer
func NewKafkaProducer(brokers []string, topic string) *KafkaProducer {
	config := sarama.NewConfig()
	config.ChannelBufferSize = 10000
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	config.Producer.RequiredAcks = sarama.WaitForLocal
	config.Producer.Compression = sarama.CompressionSnappy
	return &KafkaProducer{brokers: brokers, topic: topic, config: config}
}

// Start the kafka producer
func (conn *KafkaProducer) Start() error {
	producer, err := sarama.NewSyncProducer(conn.brokers, conn.config)
	conn.producer = producer
	// if err == nil {
	// 	go handleProducerErrors(producer)
	// }
	return err
}

// Send a new message
func (conn *KafkaProducer) Send(msg []byte) (int32, int64, error) {
	return conn.producer.SendMessage(&sarama.ProducerMessage{
		Topic: conn.topic,
		Value: sarama.ByteEncoder(msg),
	})
}

// SendMessages Send multiple messages at once to the producer resulting in better throghput for the app
func (conn *KafkaProducer) SendMessages(msgs [][]byte) error {
	messages := make([]*sarama.ProducerMessage, 0, len(msgs))
	for _, msg := range msgs {
		messages = append(messages, &sarama.ProducerMessage{
			Topic: conn.topic,
			Value: sarama.ByteEncoder(msg),
		})
	}
	return conn.producer.SendMessages(messages)
}

// Close the producer connection
func (conn *KafkaProducer) Close() error {
	if conn.producer != nil {
		return conn.producer.Close()
	}
	return nil
}

func handleProducerErrors(producer sarama.AsyncProducer) {
	for err := range producer.Errors() {
		log.Printf("Producer Error: %s\n", err.Error())
		return
	}
}
