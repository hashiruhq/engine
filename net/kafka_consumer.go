package net

import (
	"context"
	"io"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

type kafkaConsumer struct {
	brokers  []string
	topic    string
	inputs   chan kafka.Message
	consumer *kafka.Reader
}

// NewKafkaConsumer return a new Kafka consumer
func NewKafkaConsumer(brokers []string, topic string, partition int) KafkaConsumer {
	consumer := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   brokers,
		Topic:     topic,
		Partition: partition,
		MinBytes:  10,               // 10KB
		MaxBytes:  10 * 1024 * 1024, // 10MB
	})

	return &kafkaConsumer{
		brokers:  brokers,
		topic:    topic,
		consumer: consumer,
		inputs:   make(chan kafka.Message, 20000),
	}
}

func (conn *kafkaConsumer) SetOffset(offset int64) error {
	return conn.consumer.SetOffset(offset)
}

// Start the consumer
func (conn *kafkaConsumer) Start(ctx context.Context) error {
	go conn.handleMessages(ctx)
	return nil
}

// GetMessageChan returns the message channel
func (conn *kafkaConsumer) GetMessageChan() <-chan kafka.Message {
	return conn.inputs
}

// CommitMessages for the given messages
func (conn *kafkaConsumer) CommitMessages(ctx context.Context, msgs ...kafka.Message) {
	conn.consumer.CommitMessages(ctx, msgs...)
}

// Close the consumer connection
func (conn *kafkaConsumer) Close() error {
	err := conn.consumer.Close()
	close(conn.inputs)
	return err
}

/**
 * Start a supervised reader connection that automatically retries in case of an error and exists in context exit
 */
func (conn *kafkaConsumer) superviseReadingMessages(ctx context.Context) {
	for {
		err := conn.handleMessages(ctx)
		if err == io.EOF {
			// context exited or stream ended
			log.Printf("[info] [kafka] [topic:%s] Closing kafka reader: %v\n", conn.topic, err)
			break
		}
		log.Printf("[fatal] [kafka] [topic:%s] Unable to read message from reader connection with %v\n", conn.topic, err)
		break
	}
	// closing the connection
	conn.Close()
}
func (conn *kafkaConsumer) handleMessages(ctx context.Context) error {
	for {
		msg, err := conn.consumer.ReadMessage(ctx)
		// context exited or stream ended
		if err == io.EOF {
			return err
		}
		// handle any other kind of error... maybe move it in the supervisor
		if err != nil {
			kafkaErr, ok := err.(kafka.Error)
			if ok && kafkaErr.Temporary() {
				log.Printf("[warn] Unable to read message from kafka server, retrying in 1 second: %v", err)
				// wait some time before retrying
				time.Sleep(time.Second)
				continue
			} else {
				log.Fatalf("[fatal] Unable to read message from kafka server: %v", err)
				return err
			}
		}
		// send the message to the channel for processing
		conn.inputs <- msg
	}
}
