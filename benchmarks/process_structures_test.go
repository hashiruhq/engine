package benchmarks_test

import (
	"testing"

	"gitlab.com/around25/products/matching-engine/engine"
	"gitlab.com/around25/products/matching-engine/queue"

	"github.com/Shopify/sarama"
)

// @todo replace message data format to use protobuf
func BenchmarkProcessChannels(b *testing.B) {
	completed := 0
	const SIZE = 10000
	messages := make(chan []byte, SIZE)
	orders := make(chan engine.Order, SIZE)
	tradeChan := make(chan []engine.Trade, SIZE)
	done := make(chan bool)
	defer close(done)
	receiveMessages := func(messages chan<- []byte, n int) {
		for j := 0; j < n; j++ {
			messages <- []byte(`{"base": "sym","quote": "tst","id":"TST_1","price": "1312213.00010201","amount": "8483828.29993942","side": 1,"type": 1,"stop": 1,"stop_price": "13132311.00010201","funds": "101000101.33232313"}`)
		}
		close(messages)
	}
	decode := func(messages <-chan []byte, orders chan<- engine.Order) {
		for msg := range messages {
			var order engine.Order
			order.FromBinary(msg)
			orders <- order
		}
		close(orders)
	}
	processOrders := func(orders <-chan engine.Order, tradeChan chan<- []engine.Trade, n int) {
		for order := range orders {
			trades := []engine.Trade{engine.NewTrade(order.ID, order.ID, order.Amount, order.Price)}
			tradeChan <- trades
		}
		close(tradeChan)
	}
	publishTrades := func(tradeChan <-chan []engine.Trade) {
		for trades := range tradeChan {
			for _, trade := range trades {
				trade.ToBinary()
				completed++
				if completed >= b.N {
					done <- true
					return
				}
			}
		}
	}
	go receiveMessages(messages, b.N)
	go decode(messages, orders)
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
			q.Write(engine.NewEvent(&sarama.ConsumerMessage{Value: []byte(`{"base": "sym","quote": "tst","id":"TST_1","price": "1312213.00010201","amount": "8483828.29993942","side": 1,"type": 1,"stop": 1,"stop_price": "13132311.00010201","funds": "101000101.33232313"}`)}))
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
			event.SetTrades([]engine.Trade{engine.NewTrade(event.Order.ID, event.Order.ID, event.Order.Amount, event.Order.Price)})
			tradeQueue.Write(event)
		}
	}(b.N)

	go func(n int) {
		for i := 0; i < n; i++ {
			event := tradeQueue.Read()
			for _, trade := range event.Trades {
				trade.ToBinary()
			}
		}
		done <- true
	}(b.N)

	<-done
	// time.Sleep(time.Second)
}
