package net

import (
	"context"
	// "time"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/snappy"
)

// KafkaProducer structure
type kafkaProducer struct {
	producer *kafka.Writer
	brokers  []string
}

// NewKafkaProducer returns a new producer
func NewKafkaProducer(brokers []string, topic string) KafkaProducer {
	producer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:          brokers,
		Topic:            topic,
		QueueCapacity:    60000,
		BatchSize:        40000,
		// BatchTimeout:     time.Duration(100) * time.Millisecond,
		Async:            false,
		CompressionCodec: snappy.NewCompressionCodec(),
	})
	return &kafkaProducer{
		producer: producer,
		brokers:  brokers,
	}
}

// Start the kafka producer
func (conn *kafkaProducer) Start() error {
	return nil
}

// Write one or multiple messages to the topic partition
func (conn *kafkaProducer) WriteMessages(ctx context.Context, msgs ...kafka.Message) error {
	return conn.producer.WriteMessages(ctx, msgs...)
}

// Get statistics about the producer since the last time it was executed
func (conn *kafkaProducer) Stats() kafka.WriterStats {
	return conn.producer.Stats()
}

// Close the producer connection
func (conn *kafkaProducer) Close() error {
	return conn.producer.Close()
}
