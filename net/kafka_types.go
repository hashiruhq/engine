package net

import (
	"context"

	"github.com/segmentio/kafka-go"
)

// KafkaConfig godoc
type KafkaConfig struct {
	UseTLS bool              `mapstructure:"use_tls"`
	Reader KafkaReaderConfig `mapstructure:"reader"`
	Writer KafkaWriterConfig `mapstructure:"writer"`
}

// KafkaReaderConfig godoc
type KafkaReaderConfig struct {
	// The capacity of the internal message queue, defaults to 100 if none is
	// set.
	QueueCapacity int `mapstructure:"queue_capacity"` // default 100
	// Maximum amount of time to wait for new data to come when fetching batches
	// of messages from kafka.
	MaxWait int `mapstructure:"max_wait"` // default 1s
	// Min and max number of bytes to fetch from kafka in each request.
	MinBytes int `mapstructure:"min_bytes"` // 1024B = 1KB
	MaxBytes int `mapstructure:"max_bytes"` // 10485760B = 10MB
	// BackoffDelayMin optionally sets the smallest amount of time the reader will wait before
	// polling for new messages
	//
	// Default: 100ms
	ReadBackoffMin int `mapstructure:"read_backoff_min"` // default: 100 ms
	// BackoffDelayMax optionally sets the maximum amount of time the reader will wait before
	// polling for new messages
	//
	// Default: 1s
	ReadBackoffMax int `mapstructure:"read_backoff_max"` // default: 1000 ms
	// ChannelSize sets the size of the channel used to read data from the queue
	ChannelSize int `mapstructure:"channel_size"` // 20000
}

// KafkaWriterConfig godoc
type KafkaWriterConfig struct {
	// The capacity of the internal message queue, defaults to 100 if none is
	// set.
	QueueCapacity int `mapstructure:"queue_capacity"` // default 100
	// Limit on how many messages will be buffered before being sent to a
	// partition.
	//
	// The default is to use a target batch size of 100 messages.
	BatchSize int `mapstructure:"batch_size"`
	// Limit the maximum size of a request in bytes before being sent to
	// a partition.
	//
	// The default is to use a kafka default value of 1048576.
	BatchBytes int `mapstructure:"batch_bytes"`
	// Time limit on how often incomplete message batches will be flushed to
	// kafka.
	//
	// The default is to flush at least every second.
	BatchTimeout int `mapstructure:"batch_timeout"`
	// Setting this flag to true causes the WriteMessages method to never block.
	// It also means that errors are ignored since the caller will not receive
	// the returned value. Use this only if you don't care about guarantees of
	// whether the messages were written to kafka.
	Async bool `mapstructure:"async"`
}

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
