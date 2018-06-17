package trading_engine

import (
	"trading_engine/config"
	"trading_engine/net"
)

// Dispatcher holds the list of active trading markets
type Dispatcher struct {
	Markets   map[string]*TradingEngine
	Producers map[string]*net.KafkaProducer
	Consumers map[string]*net.KafkaConsumer
}

// NewDispatcher creates a new dispatcher
func NewDispatcher() Dispatcher {
	return Dispatcher{}
}

// AddMarket to the dispatcher
func (d *Dispatcher) AddMarket(base, market string) {
	pair := base + "." + market
	engine := NewTradingEngine()
	d.Markets[pair] = engine
	kafkaBroker := config.Config.Get("KAFKA_BROKER")
	kafkaOrderTopicPrefix := config.Config.Get("KAFKA_ORDER_TOPIC_PREFIX")
	kafkaOrderConsumerPrefix := config.Config.Get("KAFKA_ORDER_CONSUMER_PREFIX")
	kafkaTradeTopicPrefix := config.Config.Get("KAFKA_TRADE_TOPIC_PREFIX")

	kafkaOrderTopic := kafkaOrderTopicPrefix + base + "." + market
	kafkaTradeTopic := kafkaTradeTopicPrefix + base + "." + market
	kafkaOrderConsumer := kafkaOrderConsumerPrefix + base + "." + market

	producer := net.NewKafkaProducer([]string{kafkaBroker}, kafkaTradeTopic)
	d.Producers[pair] = producer

	consumer := net.NewKafkaConsumer([]string{kafkaBroker}, []string{kafkaOrderTopic})
	d.Consumers[pair] = consumer

	producer.Start()
	consumer.Start(kafkaOrderConsumer)
	d.processOrders(base, market)
}

func (d *Dispatcher) processOrders(base, market string) {
	pair := base + "." + market
	consumer := d.Consumers[pair]
	producer := d.Producers[pair]
	engine := d.Markets[pair]
	messageChan := consumer.GetMessageChan()
	for msg := range messageChan {
		orderMsg := net.DecodeOrder(msg.Value)
		order := NewOrder(orderMsg.ID, orderMsg.Price, orderMsg.Amount, orderMsg.Side, orderMsg.Category)
		trades := engine.Process(order)
		consumer.MarkOffset(msg, "")
		for i := range trades {
			trade := trades[i]
			tradeMsg := net.EncodeTrade(trade.MakerOrderID, trade.TakerOrderID, trade.Price, trade.Amount)
			producer.Send(tradeMsg.Encoded)
		}
	}
}
