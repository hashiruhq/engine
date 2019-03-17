package net

import (
	"log"
	"time"

	"github.com/Shopify/sarama"
	cluster "github.com/bsm/sarama-cluster"
	metrics "github.com/rcrowley/go-metrics"
)

func init() {
	metrics.UseNilMetrics = true
}

type kafkaPartitionConsumer struct {
	config   *cluster.Config
	name     string
	brokers  []string
	topics   []string
	client   *cluster.Client
	inputs   chan *sarama.ConsumerMessage
	consumer *cluster.Consumer
}

// NewKafkaPartitionConsumer return a new Kafka consumer
func NewKafkaPartitionConsumer(name string, brokers []string, topics []string) KafkaConsumer {
	config := cluster.NewConfig()
	config.Consumer.Return.Errors = true
	config.Group.Return.Notifications = true
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	// config.Net.KeepAlive = time.Duration(30) * time.Second
	config.ChannelBufferSize = 20000
	return &kafkaPartitionConsumer{
		name:    name,
		brokers: brokers,
		topics:  topics,
		config:  config,
		inputs:  make(chan *sarama.ConsumerMessage, 20000),
	}
}

// Start the consumer
func (conn *kafkaPartitionConsumer) Start() error {
	log.Printf("Starting consumer with the following config:\n - name: %s\n - brokers: %v\n - topics: %v\n", conn.name, conn.brokers, conn.topics)
	consumer, err := cluster.NewConsumer(conn.brokers, conn.name, conn.topics, conn.config)
	conn.consumer = consumer
	if err == nil {
		go conn.handleMessages()
		go conn.handleErrors()
		go conn.handleNotifications()
	} else {
		log.Println("Failed to start kafka consumer with the following error: \n", err, "\n", conn.brokers)
	}
	log.Printf("Consumer '%s' connected successfully", conn.name)

	return err
}

// GetMessageChan returns the message channel
func (conn *kafkaPartitionConsumer) GetMessageChan() <-chan *sarama.ConsumerMessage {
	return conn.inputs
}

// MarkOffset for the given message
func (conn *kafkaPartitionConsumer) MarkOffset(msg *sarama.ConsumerMessage, meta string) {
	conn.consumer.MarkOffset(msg, meta)
}

// ResetOffset will allow you to reset the Apache Kafka offset to any value you want
// Backwards or Forwards
func (conn *kafkaPartitionConsumer) ResetOffset(topic string, partition int32, offset int64, meta string) (err error) {
	consumer, err := cluster.NewConsumer(conn.brokers, conn.name, conn.topics, conn.config)
	if err != nil {
		return
	}
	log.Println("Resetting offsets")
	time.Sleep(time.Second)
	consumer.MarkPartitionOffset(topic, partition, offset, meta)
	consumer.CommitOffsets()
	consumer.ResetPartitionOffset(topic, partition, offset, meta)
	consumer.CommitOffsets()
	time.Sleep(time.Second)

	err = consumer.Close()
	if err != nil {
		return
	}
	return nil
}

// Close the consumer connection
func (conn *kafkaPartitionConsumer) Close() error {
	err := conn.consumer.Close()
	// close(conn.inputs)
	return err
}

func (conn *kafkaPartitionConsumer) handleMessages() {
	for msg := range conn.consumer.Messages() {
		conn.inputs <- msg
	}
}

func (conn *kafkaPartitionConsumer) handleErrors() {
	for err := range conn.consumer.Errors() {
		log.Printf("Error: %s\n", err.Error())
	}
}

func (conn *kafkaPartitionConsumer) handleNotifications() {
	for ntf := range conn.consumer.Notifications() {
		switch ntf.Type {
		case cluster.RebalanceStart:
			log.Printf("[info] [consumer] Kafka node is starting rebalancing...\n")
		case cluster.RebalanceOK:
			log.Printf("[info] [consumer] Rebalanced done...\n")
		case cluster.RebalanceError:
			log.Printf("[error] [consumer] Rebalancing error %+v\n", ntf)
		}
	}
}
