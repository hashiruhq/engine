package benchmarks_test

import (
	"testing"
	"trading_engine/queue"
	"trading_engine/trading_engine"

	"github.com/Shopify/sarama"
)

func BenchmarkProcessChannels(b *testing.B) {
	completed := 0
	const SIZE = 10000
	messages := make(chan []byte, SIZE)
	orders := make(chan trading_engine.Order, SIZE)
	tradeChan := make(chan []trading_engine.Trade, SIZE)
	done := make(chan bool)
	defer close(done)
	receiveMessages := func(messages chan<- []byte, n int) {
		for j := 0; j < n; j++ {
			messages <- []byte(`{"base": "sym","quote": "tst","id":"TST_1","price": "1312213.00010201","amount": "8483828.29993942","side": 1,"type": 1,"stop": 1,"stop_price": "13132311.00010201","funds": "101000101.33232313"}`)
		}
		close(messages)
	}
	jsonDecode := func(messages <-chan []byte, orders chan<- trading_engine.Order) {
		for msg := range messages {
			var order trading_engine.Order
			order.FromJSON(msg)
			orders <- order
		}
		close(orders)
	}
	processOrders := func(orders <-chan trading_engine.Order, tradeChan chan<- []trading_engine.Trade, n int) {
		for order := range orders {
			trades := []trading_engine.Trade{trading_engine.NewTrade(order.ID, order.ID, order.Amount, order.Price)}
			tradeChan <- trades
		}
		close(tradeChan)
	}
	publishTrades := func(tradeChan <-chan []trading_engine.Trade) {
		for trades := range tradeChan {
			for _, trade := range trades {
				trade.ToJSON()
				completed++
				if completed >= b.N {
					done <- true
					return
				}
			}
		}
	}
	go receiveMessages(messages, b.N)
	go jsonDecode(messages, orders)
	go processOrders(orders, tradeChan, b.N)
	go publishTrades(tradeChan)

	<-done
}

func BenchmarkProcessEventRing(b *testing.B) {
	done := make(chan bool)
	defer close(done)
	msgQueue := queue.NewBuffer(1 << 13)
	orderQueue := queue.NewBuffer(1 << 13)
	tradeQueue := queue.NewBuffer(1 << 13)

	go func(q *queue.Buffer, n int) {
		for i := 0; i < n; i++ {
			q.Write(trading_engine.NewEvent(&sarama.ConsumerMessage{Value: []byte(`{"base": "sym","quote": "tst","id":"TST_1","price": "1312213.00010201","amount": "8483828.29993942","side": 1,"type": 1,"stop": 1,"stop_price": "13132311.00010201","funds": "101000101.33232313"}`)}))
		}
	}(msgQueue, b.N)

	go func(n int) {
		for i := 0; i < n; i++ {
			event := msgQueue.Read()
			event.Decode()
			orderQueue.Write(event)
		}
	}(b.N)

	go func(n int) {
		for i := 0; i < n; i++ {
			event := orderQueue.Read()
			event.SetTrades([]trading_engine.Trade{trading_engine.NewTrade(event.Order.ID, event.Order.ID, event.Order.Amount, event.Order.Price)})
			tradeQueue.Write(event)
		}
	}(b.N)

	go func(n int) {
		for i := 0; i < n; i++ {
			event := tradeQueue.Read()
			for _, trade := range event.Trades {
				trade.ToJSON()
			}
		}
		done <- true
	}(b.N)

	<-done
	// time.Sleep(time.Second)
}
