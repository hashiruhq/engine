package server

import (
	"context"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"gitlab.com/around25/products/matching-engine/engine"
	"gitlab.com/around25/products/matching-engine/model"
	"gitlab.com/around25/products/matching-engine/net"

	"github.com/segmentio/kafka-go"
)

// MarketEngine defines how we can communicate to the trading engine for a specific market
type MarketEngine interface {
	Start(context.Context)
	Close()
	GetMessageChan() <-chan kafka.Message
	LoadMarketFromBackup() error
	Process(kafka.Message)
}

// marketEngine structure
type marketEngine struct {
	name     string
	engine   engine.TradingEngine
	inputs   chan kafka.Message
	messages chan engine.Event
	orders   chan engine.Event
	events   chan engine.Event
	backup   chan bool
	stats    chan bool
	producer net.KafkaProducer
	consumer net.KafkaConsumer
	config   MarketEngineConfig
}

// MarketEngineConfig structure
type MarketEngineConfig struct {
	producer net.KafkaProducer
	consumer net.KafkaConsumer
	config   MarketConfig
}

// NewMarketEngine open a new market
func NewMarketEngine(config MarketEngineConfig) MarketEngine {
	return &marketEngine{
		producer: config.producer,
		consumer: config.consumer,
		config:   config,
		name:     config.config.MarketID,
		engine:   engine.NewTradingEngine(config.config.MarketID, config.config.PricePrecision, config.config.VolumePrecision),
		backup:   make(chan bool),
		orders:   make(chan engine.Event, 20000),
		events:   make(chan engine.Event, 20000),
		messages: make(chan engine.Event, 20000),
	}
}

func (mkt *marketEngine) GetMessageChan() <-chan kafka.Message {
	return mkt.consumer.GetMessageChan()
}

// Start the engine for this market
func (mkt *marketEngine) Start(ctx context.Context) {
	// load last market snapshot from the backup files and update offset for the trading engine consumer
	mkt.LoadMarketFromBackup()
	if err := mkt.producer.Start(); err != nil {
		log.Fatal().Err(err).Str("section", "init:market").Str("action", "start_producer").Str("market", mkt.name).Msg("Unable to start producer")
	}
	if err := mkt.consumer.Start(ctx); err != nil {
		log.Fatal().Err(err).Str("section", "init:market").Str("action", "start_consumer").Str("market", mkt.name).Msg("Unable to start consumer")
	}
	// decode the binary value for each message received into an Order Structure
	go mkt.DecodeMessage()
	// process each order by the trading engine and forward events to the events channel
	go mkt.ProcessOrder()
	// publish events to the kafka server
	go mkt.PublishEvents()
	// start the backup scheduler
	go mkt.ScheduleBackup()
}

// Process a new message from the consumer
func (mkt *marketEngine) Process(msg kafka.Message) {
	// Monitor: Increment the number of messages that has been received by the market
	messagesQueued.WithLabelValues(mkt.name).Inc()
	mkt.messages <- engine.NewEvent(msg)
}

// Close the market by closing all communication channels
func (mkt *marketEngine) Close() {
	close(mkt.messages)
	close(mkt.orders)
	close(mkt.events)
}

// ScheduleBackup sets up an interval at which to automatically back up the market on Kafka
func (mkt *marketEngine) ScheduleBackup() {
	if mkt.config.config.Backup.Interval == 0 {
		log.Warn().Str("section", "backup").Str("action", "schedule").Str("market", mkt.name).Msg("Backup disabled for market")
		return
	}
	for {
		time.Sleep(time.Duration(mkt.config.config.Backup.Interval) * time.Minute)
		mkt.backup <- true
	}
}

// DecodeMessage decode the binary value for each message received into an Order struct
// before sending it for processing by the trading engine
//
// Message flow is unidirectional from the messages channel to the orders channel
func (mkt *marketEngine) DecodeMessage() {
	log.Debug().Str("section", "server").Str("action", "init").Str("market", mkt.name).Msg("Starting message decoder process")
	for event := range mkt.messages {
		event.Decode()
		messagesQueued.WithLabelValues(mkt.name).Dec()
		// Monitor: Increment the number of orders that are waiting to be processed
		ordersQueued.WithLabelValues(mkt.name).Inc()
		mkt.orders <- event
	}
	log.Debug().Str("section", "server").Str("action", "terminate").Str("market", mkt.name).Msg("Closed message decoder process")
}

// ProcessOrder process each order by the trading engine and forward events to the events channel
//
// Message flow is unidirectional from the orders channel to the events channel
func (mkt *marketEngine) ProcessOrder() {
	log.Debug().Str("section", "server").Str("action", "init").Str("market", mkt.name).Msg("Starting order matching process")
	var lastTopic string
	var lastPartition int32
	var lastOffset int64
	var prevOffset int64
	for {
		select {
		case <-mkt.backup:
			// Generate backup event
			if lastTopic != "" && lastOffset != prevOffset {
				market := mkt.engine.BackupMarket()
				market.Topic = lastTopic
				market.Partition = lastPartition
				market.Offset = lastOffset
				prevOffset = lastOffset
				mkt.BackupMarket(market)
				log.Debug().Str("section", "backup").Str("action", "export").Str("market", mkt.name).Msg("Snapshot created")
			}
		case event, more := <-mkt.orders:
			if !more {
				log.Debug().Str("section", "server").Str("action", "terminate").Str("market", mkt.name).Msg("Closed order matching process")
				return
			}
			order := event.Order
			if !order.Valid() {
				log.Warn().
					Str("section", "server").Str("action", "process_order").
					Str("market", mkt.name).
					Str("kafka_topic", event.Msg.Topic).
					Int("kafka_partition", event.Msg.Partition).
					Int64("kafka_offset", event.Msg.Offset).
					Dict("event", zerolog.Dict().
						Str("event_type", order.EventType.String()).
						Str("side", order.Side.String()).
						Str("type", order.Type.String()).
						Str("event_type", order.EventType.String()).
						Str("market", order.Market).
						Uint64("id", order.ID).
						Uint64("amount", order.Amount).
						Str("stop", order.Stop.String()).
						Uint64("stop_price", order.StopPrice).
						Uint64("funds", order.Funds).
						Uint64("price", order.Price),
					).
					Msg("Invalid order received, ignoring")

				// Monitor: Update order count for monitoring with prometheus
				engineOrderCount.WithLabelValues(mkt.name).Inc()
				ordersQueued.WithLabelValues(mkt.name).Dec()
				lastTopic = event.Msg.Topic
				lastPartition = int32(event.Msg.Partition)
				lastOffset = event.Msg.Offset
				continue
			}
			log.Debug().
				Str("section", "server").Str("action", "process_order").
				Str("market", mkt.name).
				Str("kafka_topic", event.Msg.Topic).
				Int("kafka_partition", event.Msg.Partition).
				Int64("kafka_offset", event.Msg.Offset).
				Dict("event", zerolog.Dict().
					Str("event_type", order.EventType.String()).
					Str("side", order.Side.String()).
					Str("type", order.Type.String()).
					Str("event_type", order.EventType.String()).
					Str("market", order.Market).
					Uint64("id", order.ID).
					Uint64("amount", order.Amount).
					Str("stop", order.Stop.String()).
					Uint64("stop_price", order.StopPrice).
					Uint64("funds", order.Funds).
					Uint64("price", order.Price),
				).
				Msg("New order")
			events := make([]model.Event, 0, 5)
			// Process each order and generate events
			mkt.engine.ProcessEvent(event.Order, &events)
			event.SetEvents(events)
			// Monitor: Update order count for monitoring with prometheus
			engineOrderCount.WithLabelValues(mkt.name).Inc()
			ordersQueued.WithLabelValues(mkt.name).Dec()
			eventsQueued.WithLabelValues(mkt.name).Add(float64(len(event.Events)))
			// send generated events for storage
			mkt.events <- event
			lastTopic = event.Msg.Topic
			lastPartition = int32(event.Msg.Partition)
			lastOffset = event.Msg.Offset
		}
	}
}

// PublishEvents listens for new events from the trading engine and publishes them to the Kafka server
func (mkt *marketEngine) PublishEvents() {
	log.Debug().Str("section", "server").Str("action", "init").Str("market", mkt.name).Msg("Starting event publisher process")
	var lastAskID uint64
	var lastBidID uint64
	for event := range mkt.events {
		events := make([]kafka.Message, len(event.Events))
		for index, ev := range event.Events {
			logEvent := zerolog.Dict()
			switch ev.Type {
			case model.EventType_OrderStatusChange:
				{
					payload := ev.GetOrderStatus()
					logEvent = logEvent.
						Uint64("id", payload.ID).
						Uint64("owner_id", payload.OwnerID).
						Str("type", payload.Type.String()).
						Str("side", payload.Side.String()).
						Str("status", payload.Status.String()).
						Uint64("price", payload.Price).
						Uint64("funds", payload.Funds).
						Uint64("amount", payload.Amount)
				}
			case model.EventType_OrderActivated:
				{
					payload := ev.GetOrderActivation()
					logEvent = logEvent.
						Uint64("id", payload.ID).
						Uint64("owner_id", payload.OwnerID).
						Str("type", payload.Type.String()).
						Str("side", payload.Side.String()).
						Str("status", payload.Status.String()).
						Uint64("price", payload.Price).
						Uint64("funds", payload.Funds).
						Uint64("amount", payload.Amount)
				}
			case model.EventType_NewTrade:
				{
					trade := ev.GetTrade()
					logEvent = logEvent.
						Uint64("seqid", trade.SeqID).
						Str("taker_side", trade.TakerSide.String()).
						Uint64("ask_id", trade.AskID).
						Uint64("ask_owner_id", trade.AskOwnerID).
						Uint64("bid_id", trade.BidID).
						Uint64("bid_owner_id", trade.BidOwnerID).
						Uint64("price", trade.Price).
						Uint64("amount", trade.Amount)
					if lastAskID == trade.AskID && lastBidID == trade.BidID {
						log.Error().Str("section", "engine").Str("action", "post:trade:check").
							Str("market", mkt.name).
							Uint64("last_order_id", event.Order.ID).
							Str("kafka_topic", event.Msg.Topic).
							Int("kafka_partition", event.Msg.Partition).
							Int64("kafka_offset", event.Msg.Offset).
							Str("event_type", ev.Type.String()).
							Uint64("event_seqid", ev.SeqID).
							Int64("event_timestamp", ev.CreatedAt).
							Dict("event", logEvent).
							Msg("An bid order matched with the same sell order twice. Orderbook is in an inconsistent state.")
					}
					lastBidID = trade.BidID
					lastAskID = trade.AskID
				}
			}
			log.Debug().Str("section", "server").Str("action", "publish").
				Str("market", mkt.name).
				Uint64("last_order_id", event.Order.ID).
				Str("kafka_topic", event.Msg.Topic).
				Int("kafka_partition", event.Msg.Partition).
				Int64("kafka_offset", event.Msg.Offset).
				Str("event_type", ev.Type.String()).
				Uint64("event_seqid", ev.SeqID).
				Int64("event_timestamp", ev.CreatedAt).
				Dict("event", logEvent).
				Msg("Generated event")
			rawTrade, _ := ev.ToBinary() // @todo add better error handling on encoding
			events[index] = kafka.Message{
				Value: rawTrade,
			}
		}
		err := mkt.producer.WriteMessages(context.Background(), events...)
		if err != nil {
			log.Fatal().Err(err).Str("section", "server").Str("action", "publish").Str("market", mkt.name).Msg("Unable to publish events")
		}

		// Monitor: Update the number of events processed after sending them back to Kafka
		eventCount := float64(len(event.Events))
		eventsQueued.WithLabelValues(mkt.name).Sub(eventCount)
		engineEventCount.WithLabelValues(mkt.name).Add(eventCount)
	}
	log.Info().Str("section", "server").Str("action", "terminate").Str("market", mkt.name).Msg("Closing event publisher process")
}
