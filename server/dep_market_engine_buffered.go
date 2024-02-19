package server

/**
THIS FILE IS DEPRECATED AND SHOULD NOT BE USED EXCEPT FOR TESTING
*/

// import (
// 	"io/ioutil"
// 	"log"
// 	"os"
// 	"time"

// 	"github.com/hashiruhq/engine/engine"
// 	"github.com/hashiruhq/engine/net"
// 	"github.com/hashiruhq/engine/queue"

// 	"github.com/Shopify/sarama"
// )

// // bufferedMarketEngine structure
// type bufferedMarketEngine struct {
// 	name           string
// 	engine         engine.TradingEngine
// 	messages       *queue.Buffer
// 	orders         *queue.Buffer
// 	trades         *queue.Buffer
// 	backup         chan bool
// 	producer       net.KafkaProducer
// 	backupProducer net.KafkaProducer
// 	consumer       net.KafkaConsumer
// 	config         MarketEngineConfig
// }

// // NewBufferedMarketEngine open a new market
// func NewBufferedMarketEngine(config MarketEngineConfig) MarketEngine {
// 	return &bufferedMarketEngine{
// 		producer: config.producer,
// 		consumer: config.consumer,
// 		config:   config,
// 		name:     config.config.Base + "_" + config.config.Quote,
// 		engine:   engine.NewTradingEngine(),
// 		backup:   make(chan bool),
// 		messages: queue.NewBuffer(1 << 15),
// 		orders:   queue.NewBuffer(1 << 15),
// 		trades:   queue.NewBuffer(1 << 15),
// 	}
// }

// // Start the engine for this market
// func (mkt *bufferedMarketEngine) Start() {
// 	// listen for producer errors when publishing trades
// 	go mkt.ReceiveProducerErrors()
// 	// decode the binary value for each message received into an Order Structure
// 	go mkt.DecodeMessage()
// 	// process each order by the trading engine and forward trades to the trades channel
// 	go mkt.ProcessOrder()
// 	// publish trades to the kafka server
// 	go mkt.PublishTrades()
// 	// start the backup scheduler
// 	go mkt.ScheduleBackup()
// }

// // Process a new message from the consumer
// func (mkt *bufferedMarketEngine) Process(msg *sarama.ConsumerMessage) {
// 	// Monitor: Increment the number of messages that has been received by the market
// 	messagesQueued.WithLabelValues(mkt.name).Inc()
// 	mkt.messages.Write(engine.NewEvent(msg))
// }

// // Close the market by closing all communication channels
// func (mkt *bufferedMarketEngine) Close() {
// 	// do nothing
// }

// // ReceiveProducerErrors listens for error messages trades that have not been published
// // and sends them on another queue to processing later.
// //
// // @todo Optionally we can send them to another service, like an in memory database (Redis) so
// // that we don't retry to push then to a broken Kafka instance
// func (mkt *bufferedMarketEngine) ReceiveProducerErrors() {
// 	errors := mkt.producer.Errors()
// 	for err := range errors {
// 		log.Printf("Error sending trade on '%s' market: %v", mkt.name, err)
// 	}
// }

// // ScheduleBackup sets up an interval at which to automatically back up the market on Kafka
// func (mkt *bufferedMarketEngine) ScheduleBackup() {
// 	if mkt.config.config.Backup.Interval == 0 {
// 		log.Printf("Backup disabled for '%s' market. Please set backup interval to a positive value.", mkt.name)
// 		return
// 	}
// 	for {
// 		time.Sleep(time.Duration(mkt.config.config.Backup.Interval) * time.Minute)
// 		mkt.backup <- true
// 	}
// }

// // DecodeMessage decode the binary value for each message received into an Order struct
// // before sending it for processing by the trading engine
// //
// // Message flow is unidirectional from the messages channel to the orders channel
// func (mkt *bufferedMarketEngine) DecodeMessage() {
// 	for {
// 		event := mkt.messages.Read()
// 		event.Decode()
// 		messagesQueued.WithLabelValues(mkt.name).Dec()
// 		// Monitor: Increment the number of orders that are waiting to be processed
// 		ordersQueued.WithLabelValues(mkt.name).Inc()
// 		mkt.orders.Write(event)
// 	}
// }

// // ProcessOrder process each order by the trading engine and forward trades to the trades channel
// //
// // Message flow is unidirectional from the orders channel to the trades channel
// func (mkt *bufferedMarketEngine) ProcessOrder() {
// 	var lastTopic string
// 	var lastPartition int32
// 	var lastOffset int64
// 	var prevOffset int64
// 	for {
// 		select {
// 		case <-mkt.backup:
// 			// Generate backup event
// 			if lastTopic != "" && lastOffset != prevOffset {
// 				market := mkt.engine.BackupMarket()
// 				market.Topic = lastTopic
// 				market.Partition = lastPartition
// 				market.Offset = lastOffset
// 				prevOffset = lastOffset
// 				mkt.BackupMarket(market)
// 				log.Printf("[Backup] [%s] Snapshot created\n", mkt.name)
// 			} else {
// 				log.Printf("[Backup] [%s] Skipped. No changes since last backup\n", mkt.name)
// 			}
// 		default:
// 			event := mkt.orders.Read()
// 			lastTopic = event.Msg.Topic
// 			lastPartition = event.Msg.Partition
// 			lastOffset = event.Msg.Offset
// 			// Process each order and generate trades
// 			// event.SetTrades(mkt.engine.Process(event.Order))
// 			// Monitor: Update order count for monitoring with prometheus
// 			engineOrderCount.WithLabelValues(mkt.name).Inc()
// 			ordersQueued.WithLabelValues(mkt.name).Dec()
// 			tradesQueued.WithLabelValues(mkt.name).Add(float64(len(event.Trades)))
// 			// send trades for storage
// 			mkt.trades.Write(event)
// 		}
// 	}
// }

// // PublishTrades listens for new trades from the trading engine and publishes them to the Kafka server
// func (mkt *bufferedMarketEngine) PublishTrades() {
// 	for {
// 		event := mkt.trades.Read()
// 		for _, trade := range event.Trades {
// 			rawTrade, _ := trade.ToBinary() // @todo threat error on encoding object
// 			mkt.producer.Input() <- &sarama.ProducerMessage{
// 				Topic: mkt.config.config.Publish.Topic,
// 				Value: sarama.ByteEncoder(rawTrade),
// 			}
// 		}
// 		// mark the order as completed only after the trades have been sent to the kafka server
// 		mkt.consumer.MarkOffset(event.Msg, "")

// 		// Monitor: Update the number of trades processed after sending them back to Kafka
// 		tradeCount := float64(len(event.Trades))
// 		tradesQueued.WithLabelValues(mkt.name).Sub(tradeCount)
// 		engineTradeCount.WithLabelValues(mkt.name).Add(tradeCount)
// 	}
// }

// // BackupMarket saves the given snapshot of the order book as binary into the backups folder with the name of the market pair
// // - It first saves into a temporary file before moving the file to the final localtion
// func (mkt *bufferedMarketEngine) BackupMarket(market engine.MarketBackup) error {
// 	file := mkt.config.config.Backup.Path + ".tmp"
// 	rawMarket, _ := market.ToBinary()
// 	ioutil.WriteFile(file, rawMarket, 0644)
// 	os.Rename(mkt.config.config.Backup.Path+".tmp", mkt.config.config.Backup.Path)
// 	return nil
// }

// // LoadMarketFromBackup from a backup file and update the order book with the given data
// // - Also reset the Kafka partition offset to the one from the backup and replay the orders to recreate the latest state
// func (mkt *bufferedMarketEngine) LoadMarketFromBackup() (err error) {
// 	log.Printf("Loading market %s from backup file: %s", mkt.name, mkt.config.config.Backup.Path)
// 	file := mkt.config.config.Backup.Path
// 	content, err := ioutil.ReadFile(file)
// 	if err != nil {
// 		return
// 	}

// 	var market engine.MarketBackup
// 	market.FromBinary(content)

// 	// load all records from the backup into the order book
// 	mkt.LoadMarket(market)

// 	// mark the last message that has been processed by the engine to the one saved in the backup file
// 	err = mkt.consumer.ResetOffset(market.Topic, market.Partition, market.Offset, "")
// 	if err != nil {
// 		log.Printf("Unable to reset offset for the '%s' market on '%s' topic and partition '%d' to offset '%d'", mkt.name, market.Topic, market.Partition, market.Offset)
// 		return
// 	}

// 	log.Printf("Market %s loaded from backup", mkt.name)

// 	return nil
// }

// func (mkt *bufferedMarketEngine) LoadMarket(market engine.MarketBackup) {
// 	mkt.engine.LoadMarket(market)
// }
