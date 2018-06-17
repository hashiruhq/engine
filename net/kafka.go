package net

import (
	"encoding/json"
	"log"

	"github.com/Shopify/sarama"

	cluster "github.com/bsm/sarama-cluster"
)

// OrderMessage structure
type OrderMessage struct {
	ID       string  `json:"id"`
	Price    float64 `json:"price"`
	Amount   float64 `json:"amount"`
	Side     int     `json:"side"`
	Category int     `json:"category"`
	Encoded  []byte
}

// EncodeOrder an order as a message
func EncodeOrder(id string, price float64, amount float64, side int) *OrderMessage {
	msg := OrderMessage{ID: id, Price: price, Amount: amount, Side: side}
	encoded, _ := json.Marshal(msg)
	msg.Encoded = encoded
	return &msg
}

// DecodeOrder byte array
func DecodeOrder(encoded []byte) *OrderMessage {
	msg := OrderMessage{Encoded: encoded}
	json.Unmarshal(encoded, &msg)
	return &msg
}

// Length returns the length of the message
func (msg *OrderMessage) Length() int {
	return len(msg.Encoded)
}

// Encode the message and return the byte array
func (msg *OrderMessage) Encode() ([]byte, error) {
	return msg.Encoded, nil
}

// TradeMessage structure
type TradeMessage struct {
	MakerOrderID string  `json:"maker_order_id"`
	TakerOrderID string  `json:"taker_order_id"`
	Price        float64 `json:"price"`
	Amount       float64 `json:"amount"`
	Encoded      []byte
}

// EncodeTrade an trade as a message
func EncodeTrade(maker string, taker string, price float64, amount float64) *TradeMessage {
	msg := TradeMessage{MakerOrderID: maker, TakerOrderID: taker, Price: price, Amount: amount}
	encoded, _ := json.Marshal(msg)
	msg.Encoded = encoded
	return &msg
}

// DecodeTrade byte array
func DecodeTrade(encoded []byte) *TradeMessage {
	msg := TradeMessage{Encoded: encoded}
	json.Unmarshal(encoded, &msg)
	return &msg
}

// Length returns the length of the message
func (msg *TradeMessage) Length() int {
	return len(msg.Encoded)
}

// Encode the message and return the byte array
func (msg *TradeMessage) Encode() ([]byte, error) {
	return msg.Encoded, nil
}

// KafkaProducer structure
type KafkaProducer struct {
	brokers  []string
	topic    string
	producer sarama.SyncProducer
	config   *sarama.Config
}

// NewKafkaProducer returns a new producer
func NewKafkaProducer(brokers []string, topic string) *KafkaProducer {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForLocal
	return &KafkaProducer{brokers: brokers, topic: topic, config: config}
}

// Start the kafka producer
func (conn *KafkaProducer) Start() error {
	producer, err := sarama.NewSyncProducer(conn.brokers, conn.config)
	conn.producer = producer
	// if err == nil {
	// 	go handleProducerErrors(producer)
	// }
	return err
}

// Send a new message
func (conn *KafkaProducer) Send(msg []byte) (int32, int64, error) {
	return conn.producer.SendMessage(&sarama.ProducerMessage{
		Topic: conn.topic,
		Value: sarama.ByteEncoder(msg),
	})
}

// Close the producer connection
func (conn *KafkaProducer) Close() error {
	if conn.producer != nil {
		return conn.producer.Close()
	}
	return nil
}

func handleProducerErrors(producer sarama.AsyncProducer) {
	for err := range producer.Errors() {
		log.Printf("Producer Error: %s\n", err.Error())
		return
	}
}

// KafkaConsumer structure
type KafkaConsumer struct {
	config   *cluster.Config
	brokers  []string
	topics   []string
	consumer *cluster.Consumer
}

// NewKafkaConsumer return a new Kafka consumer
func NewKafkaConsumer(brokers []string, topics []string) *KafkaConsumer {
	config := cluster.NewConfig()
	config.Consumer.Return.Errors = false
	config.Group.Return.Notifications = false
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	return &KafkaConsumer{brokers: brokers, topics: topics, config: config}
}

// Start the consumer
func (conn *KafkaConsumer) Start(name string) error {
	consumer, err := cluster.NewConsumer(conn.brokers, name, conn.topics, conn.config)
	conn.consumer = consumer
	if err == nil {
		go handleErrors(consumer)
		go handleNotifications(consumer)
	}
	return err
}

// GetMessageChan returns the message channel
func (conn *KafkaConsumer) GetMessageChan() <-chan *sarama.ConsumerMessage {
	return conn.consumer.Messages()
}

// MarkOffset for the given message
func (conn *KafkaConsumer) MarkOffset(msg *sarama.ConsumerMessage, meta string) {
	conn.consumer.MarkOffset(msg, meta)
}

// Close the consumer connection
func (conn *KafkaConsumer) Close() error {
	return conn.consumer.Close()
}

func handleErrors(consumer *cluster.Consumer) {
	for err := range consumer.Errors() {
		log.Printf("Error: %s\n", err.Error())
	}
}

func handleNotifications(consumer *cluster.Consumer) {
	for ntf := range consumer.Notifications() {
		log.Printf("Rebalanced: %+v\n", ntf)
	}
}
