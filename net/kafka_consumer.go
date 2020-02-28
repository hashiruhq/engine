package net

import (
	"context"
	"io"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"
)

type kafkaConsumer struct {
	brokers  []string
	topic    string
	inputs   chan kafka.Message
	consumer *kafka.Reader
}

// NewKafkaConsumer return a new Kafka consumer
func NewKafkaConsumer(cfg KafkaReaderConfig, brokers []string, topic string, partition int) KafkaConsumer {
	consumer := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        brokers,
		Topic:          topic,
		Partition:      partition,
		MaxWait:        time.Duration(cfg.MaxWait) * time.Millisecond,
		QueueCapacity:  cfg.QueueCapacity,
		MinBytes:       cfg.MinBytes,
		MaxBytes:       cfg.MaxBytes,
		ReadBackoffMin: time.Duration(cfg.ReadBackoffMin) * time.Millisecond,
		ReadBackoffMax: time.Duration(cfg.ReadBackoffMax) * time.Millisecond,
	})

	return &kafkaConsumer{
		brokers:  brokers,
		topic:    topic,
		consumer: consumer,
		inputs:   make(chan kafka.Message, cfg.ChannelSize),
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
			log.Info().
				Err(err).
				Str("section", "kafka").
				Str("topic", conn.topic).
				Msg("Context exited of message stream ended. Closing consumer")
			break
		}
		log.Error().
			Err(err).
			Str("section", "kafka").
			Str("topic", conn.topic).
			Msg("Unable to read message from reader connection. Closing consumer")
		break
	}
	// closing the connection
	conn.Close()
	close(conn.inputs)
}

func (conn *kafkaConsumer) handleMessages(ctx context.Context) error {
	log.Info().
		Str("section", "kafka").
		Str("topic", conn.topic).
		Int64("offset", conn.consumer.Offset()).
		Msg("Starting message consumer")
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
				log.Warn().
					Err(kafkaErr).
					Str("section", "kafka").
					Str("topic", conn.topic).
					Bool("temp", true).
					Msg("Unable to read message from reader connection. Retrying in 1 second")
				// wait some time before retrying
				time.Sleep(time.Second)
				continue
			} else {
				log.Error().
					Err(err).
					Str("section", "kafka").
					Str("topic", conn.topic).
					Bool("temp", false).
					Msg("Unable to read message from kafka server. Exiting message reader loop.")
				return err
			}
		}
		// send the message to the channel for processing
		conn.inputs <- msg
	}
}
