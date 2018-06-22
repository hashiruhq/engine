package net

// import (
// 	"github.com/Shopify/sarama"
// )

// // KafkaProducer structure
// type KafkaSyncProducer struct {
// 	brokers  []string
// 	topic    string
// 	input    chan []byte
// 	errors   chan *sarama.ProducerError
// 	producer sarama.SyncProducer
// 	config   *sarama.Config
// }

// // NewKafkaSyncProducer returns a new producer
// func NewKafkaSyncProducer(brokers []string, topic string) KafkaProducer {
// 	config := sarama.NewConfig()
// 	config.ChannelBufferSize = 10000
// 	config.Producer.Return.Successes = true
// 	config.Producer.Return.Errors = true
// 	config.Producer.RequiredAcks = sarama.WaitForAll
// 	config.Producer.Compression = sarama.CompressionSnappy
// 	return &KafkaSyncProducer{
// 		brokers: brokers,
// 		topic:   topic,
// 		config:  config,
// 		input:   make(chan []byte),
// 		errors:  make(chan *sarama.ProducerError),
// 	}
// }

// // Start the kafka producer
// func (conn *KafkaSyncProducer) Start() error {
// 	producer, err := sarama.NewSyncProducer(conn.brokers, conn.config)
// 	conn.producer = producer
// 	go conn.dispatch()
// 	return err
// }

// // Input a new message to the producer
// func (conn *KafkaSyncProducer) Input() chan<- []byte {
// 	return conn.input
// }

// // Errors returns the error channel
// func (conn *KafkaSyncProducer) Errors() <-chan *sarama.ProducerError {
// 	return conn.errors
// }

// // Close the producer connection
// func (conn *KafkaSyncProducer) Close() error {
// 	if conn.producer != nil {
// 		err := conn.producer.Close()
// 		close(conn.input)
// 		return err
// 	}
// 	return nil
// }

// func (conn *KafkaSyncProducer) dispatch() {
// 	for msg := range conn.input {
// 		producerMsg := &sarama.ProducerMessage{
// 			Topic: conn.topic,
// 			Value: sarama.ByteEncoder(msg),
// 		}
// 		partition, offset, err := conn.producer.SendMessage(producerMsg)
// 		if err != nil {
// 			conn.errors <- err.(*sarama.ProducerError)
// 		}
// 	}
// }
