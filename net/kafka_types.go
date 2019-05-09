package net

import (
	"context"

	"github.com/segmentio/kafka-go"
)

// KafkaProducer inferface
type KafkaProducer interface {
	Start() error
	WriteMessages(context.Context, ...kafka.Message) error
	Close() error
}

// KafkaConsumer interface
type KafkaConsumer interface {
	Start(ctx context.Context) error
	SetOffset(offset int64) error
	GetMessageChan() <-chan kafka.Message
	CommitMessages(context.Context, ...kafka.Message)
	Close() error
}
