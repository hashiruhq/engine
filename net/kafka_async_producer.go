package net

import (
	"github.com/Shopify/sarama"
)

// KafkaAsyncProducer structure
type KafkaAsyncProducer struct {
	brokers  []string
	topic    string
	input    chan []byte
	producer sarama.AsyncProducer
	config   *sarama.Config
}

// NewKafkaAsyncProducer returns a new producer
func NewKafkaAsyncProducer(brokers []string, topic string) KafkaProducer {
	config := sarama.NewConfig()
	// config.ChannelBufferSize = 10000
	// config.Producer.Return.Successes = true
	config.Producer.Retry.Max = 5
	config.Producer.Return.Errors = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	// config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Compression = sarama.CompressionSnappy
	producer := &KafkaAsyncProducer{
		brokers: brokers,
		topic:   topic,
		config:  config,
		input:   make(chan []byte),
	}

	return producer
}

// Start the kafka producer
func (conn *KafkaAsyncProducer) Start() error {
	producer, err := sarama.NewAsyncProducer(conn.brokers, conn.config)
	conn.producer = producer
	go conn.dispatch()
	return err
}

// Input a new message to the producer
func (conn *KafkaAsyncProducer) Input() chan<- []byte {
	return conn.input
}

// Errors returns the error channel
func (conn *KafkaAsyncProducer) Errors() <-chan *sarama.ProducerError {
	return conn.producer.Errors()
}

// Close the producer connection
func (conn *KafkaAsyncProducer) Close() error {
	if conn.producer != nil {
		err := conn.producer.Close()
		close(conn.input)
		return err
	}
	return nil
}

func (conn *KafkaAsyncProducer) dispatch() {
	for msg := range conn.input {
		conn.producer.Input() <- &sarama.ProducerMessage{
			Topic: conn.topic,
			Value: sarama.ByteEncoder(msg),
		}
	}
}
