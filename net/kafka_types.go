package net

import "github.com/Shopify/sarama"

// KafkaProducer inferface
type KafkaProducer interface {
	Start() error
	Input() chan<- *sarama.ProducerMessage
	Errors() <-chan *sarama.ProducerError
	Close() error
}

// KafkaConsumer interface
type KafkaConsumer interface {
	Start() error
	GetMessageChan() <-chan *sarama.ConsumerMessage
	MarkOffset(msg *sarama.ConsumerMessage, meta string)
	ResetOffset(topic string, partition int32, offset int64, meta string) error
	Close() error
}
