package net

import "github.com/Shopify/sarama"

// KafkaProducer inferface
type KafkaProducer interface {
	Start() error
	Input() chan<- []byte
	Errors() <-chan *sarama.ProducerError
	Close() error
}

// KafkaConsumer interface
type KafkaConsumer interface {
	Start(name string) error
	GetMessageChan() <-chan *sarama.ConsumerMessage
	MarkOffset(msg *sarama.ConsumerMessage, meta string)
	Close() error
}
