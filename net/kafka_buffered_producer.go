package net

import (
	"sync"
)

const MAX_BUFFER_SIZE = 10000

// KafkaBufferedProducer structure
// @todo Account for any errors when sending to the Kafka server
// @todo Check if sending the records in parallel has any affect on data consistency
type KafkaBufferedProducer struct {
	brokers    []string
	topic      string
	buffer     [][]byte
	inProgress sync.WaitGroup
	lock       sync.Mutex
	producer   *KafkaProducer
}

// NewKafkaBufferedProducer returns a new producer
func NewKafkaBufferedProducer(brokers []string, topic string) *KafkaBufferedProducer {
	producer := NewKafkaProducer(brokers, topic)
	buffer := make([][]byte, 0, MAX_BUFFER_SIZE)
	return &KafkaBufferedProducer{brokers: brokers, topic: topic, producer: producer, buffer: buffer}
}

// Start the kafka producer
func (conn *KafkaBufferedProducer) Start() error {
	return conn.producer.Start()
}

// SendMessages Send multiple messages at once to the producer resulting in better throughput for the app
func (conn *KafkaBufferedProducer) SendMessages(msgs [][]byte) {
	if len(msgs)+len(conn.buffer) <= MAX_BUFFER_SIZE {
		conn.buffer = append(conn.buffer, msgs...)
		return
	}
	conn.sendMessages()
	if len(msgs) <= MAX_BUFFER_SIZE {
		conn.buffer = append(conn.buffer, msgs...)
		return
	}
	conn.buffer = append(conn.buffer, msgs[0:MAX_BUFFER_SIZE]...)
	conn.SendMessages(msgs[MAX_BUFFER_SIZE:])
}

func (conn *KafkaBufferedProducer) sendMessages() {
	if len(conn.buffer) == 0 {
		return
	}
	conn.inProgress.Add(1)
	go func(buffer [][]byte) {
		conn.lock.Lock()
		conn.producer.SendMessages(buffer)
		conn.lock.Unlock()
		conn.inProgress.Done()
	}(conn.buffer)
	conn.buffer = make([][]byte, 0, MAX_BUFFER_SIZE)
}

// Flush the remaining messages from the buffer
func (conn *KafkaBufferedProducer) Flush() {
	conn.sendMessages()
}

// Wait for all producers to finish sending their message payload
func (conn *KafkaBufferedProducer) Wait() {
	conn.inProgress.Wait()
}

// Close the producer connection
func (conn *KafkaBufferedProducer) Close() error {
	if conn.producer != nil {
		return conn.producer.Close()
	}
	return nil
}
