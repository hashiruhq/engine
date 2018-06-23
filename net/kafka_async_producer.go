package net

import (
	"github.com/Shopify/sarama"
)

// KafkaAsyncProducer structure
type kafkaAsyncProducer struct {
	input    chan []byte
	topic    string
	producer sarama.AsyncProducer
	config   *sarama.Config
	brokers  []string
}

// NewKafkaAsyncProducer returns a new producer
func NewKafkaAsyncProducer(brokers []string, topic string) KafkaProducer {
	config := sarama.NewConfig()
	// config.ChannelBufferSize = 10000
	config.Producer.Return.Successes = false
	config.Producer.Retry.Max = 5
	config.Producer.Return.Errors = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRoundRobinPartitioner
	config.Producer.Compression = sarama.CompressionSnappy
	return &kafkaAsyncProducer{
		brokers: brokers,
		topic:   topic,
		config:  config,
		input:   make(chan []byte),
	}
}

// Start the kafka producer
func (conn *kafkaAsyncProducer) Start() error {
	producer, err := sarama.NewAsyncProducer(conn.brokers, conn.config)
	conn.producer = producer
	go conn.dispatch()
	return err
}

// Input a new message to the producer
func (conn *kafkaAsyncProducer) Input() chan<- []byte {
	return conn.input
}

// Errors returns the error channel
func (conn *kafkaAsyncProducer) Errors() <-chan *sarama.ProducerError {
	return conn.producer.Errors()
}

// Close the producer connection
func (conn *kafkaAsyncProducer) Close() error {
	if conn.producer != nil {
		err := conn.producer.Close()
		close(conn.input)
		return err
	}
	return nil
}

func (conn *kafkaAsyncProducer) dispatch() {
	for msg := range conn.input {
		conn.producer.Input() <- &sarama.ProducerMessage{
			Topic: conn.topic,
			Value: sarama.ByteEncoder(msg),
		}
	}
}