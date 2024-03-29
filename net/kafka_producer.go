package net

import (
	"context"
	"crypto/tls"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/snappy"
)

// KafkaProducer structure
type kafkaProducer struct {
	producer *kafka.Writer
	brokers  []string
}

// NewKafkaProducer returns a new producer
func NewKafkaProducer(cfg KafkaWriterConfig, brokers []string, useTLS bool, topic string) KafkaProducer {
	var dialer *kafka.Dialer
	if useTLS {
		tlsCfg := &tls.Config{
			//InsecureSkipVerify: true,
		}
		dialer = &kafka.Dialer{
			Timeout:   10 * time.Second,
			DualStack: true,
			TLS:       tlsCfg,
		}
	}
	producer := kafka.NewWriter(kafka.WriterConfig{
		Dialer:           dialer,
		Brokers:          brokers,
		Topic:            topic,
		QueueCapacity:    cfg.QueueCapacity,
		BatchSize:        cfg.BatchSize,
		BatchTimeout:     time.Duration(cfg.BatchTimeout) * time.Millisecond,
		Async:            cfg.Async,
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
