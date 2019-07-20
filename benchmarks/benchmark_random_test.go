package benchmarks_test

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/golang/protobuf/proto"
	"gitlab.com/around25/products/matching-engine/engine"
	"gitlab.com/around25/products/matching-engine/model"
)

var arr []model.Order = make([]model.Order, 0, 2000000)
var msgs [][]byte = make([][]byte, 0, 2000000)
var ngin = engine.NewTradingEngine("btcusd", 8, 8)

func init() {
	GenerateRandomRecordsInFile(BENCHMARK_TEST_FILE, KAFKA_CONSUMER_MARKET, 2000000)
	file, err := filepath.Abs(BENCHMARK_TEST_FILE)
	if err != nil {
		panic(err.Error())
	}
	fh, err := os.Open(file)
	if err != nil {
		panic(err.Error())
	}
	defer fh.Close()
	bf := bufio.NewReader(fh)
	for j := 0; j < 2000000; j++ {
		msg, _, err := bf.ReadLine()
		if err != nil {
			log.Fatalln(err)
		}
		data := make([]byte, 1000)
		base64.StdEncoding.Decode(data, msg)
		order := model.Order{}
		order.FromBinary(data)
		data, _ = order.ToBinary()
		msgs = append(msgs, data)
		arr = append(arr, order)
	}
}

func BenchmarkDecodeFromProto(benchmark *testing.B) {
	order := &model.Order{
		ID:        1,
		Market:    "btc-usd",
		Amount:    848382829993942,
		Price:     131221300010201,
		Side:      model.MarketSide_Buy,
		Type:      model.OrderType_Limit,
		EventType: model.CommandType_NewOrder,
		Stop:      model.StopLoss_Loss,
		StopPrice: 1313231100010201,
		Funds:     10100010133232313,
	}
	binary, _ := proto.Marshal(order)
	for i := 0; i < benchmark.N; i++ {
		proto.Unmarshal(binary, order)
	}
}

func BenchmarkEncodeToProto(benchmark *testing.B) {
	order := &model.Order{
		ID:        1,
		Market:    "btc-usd",
		Amount:    848382829993942,
		Price:     131221300010201,
		Side:      model.MarketSide_Buy,
		Type:      model.OrderType_Limit,
		EventType: model.CommandType_NewOrder,
		Stop:      model.StopLoss_Loss,
		StopPrice: 1313231100010201,
		Funds:     10100010133232313,
	}
	for i := 0; i < benchmark.N; i++ {
		proto.Marshal(order)
	}
}

func BenchmarkWithRandomData(benchmark *testing.B) {
	startTime := time.Now().UnixNano()
	events := make([]model.Event, 0, 100)
	processing_events := make([]model.Event, 0, 100)
	// ngin := engine.NewTradingEngine()
	if benchmark.N > 2000000 {
		panic("Need more data to test with")
	}
	for j := 0; j < benchmark.N; j++ {
		ngin.Process(arr[j], &events)
		if len(events) >= cap(events)/2+1 {
			copy(processing_events, events)
			events = events[0:0]
			processing_events = processing_events[0:0]
		}
	}
	PrintOrderLogs(ngin, benchmark.N, startTime)
}

func BenchmarkWithDecodeAndEncodeRandomData(benchmark *testing.B) {
	startTime := time.Now().UnixNano()
	events := make([]model.Event, 0, 200)
	processing_events := make([]model.Event, 0, 200)
	// ngin := engine.NewTradingEngine()
	for j := 0; j < benchmark.N; j++ {
		order := model.Order{}
		order.FromBinary(msgs[j])
		ngin.Process(order, &events)
		if len(events) >= cap(events)/2+1 {
			for _, event := range events {
				event.ToBinary()
			}
			copy(processing_events, events)
			events = events[0:0]
			processing_events = processing_events[0:0]
		}
	}
	PrintOrderLogs(ngin, benchmark.N, startTime)
}

func BenchmarkTimestamp(benchmark *testing.B) {
	for j := 0; j < benchmark.N; j++ {
		time.Now().UTC().Unix()
	}
}

func PrintOrderLogs(engine engine.TradingEngine, ordersCompleted int, startTime int64) {
	endTime := time.Now().UnixNano()
	timeout := (float64)(float64(time.Nanosecond) * float64(endTime-startTime) / float64(time.Second))
	fmt.Printf(
		"Total Orders: %d\n"+
			// "Total Trades: %d\n"+
			"Orders/second: %f\n"+
			// "Trades/second: %f\n"+
			"Pending Buy: %d\n"+
			"Lowest Ask: %d\n"+
			"Pending Sell: %d\n"+
			"Highest Bid: %d\n"+
			"Duration (seconds): %f\n\n",
		ordersCompleted,
		// tradesCompleted,
		float64(ordersCompleted)/timeout,
		// float64(tradesCompleted)/timeout,
		engine.GetOrderBook().GetMarket()[0].Len(),
		engine.GetOrderBook().GetLowestAsk(),
		engine.GetOrderBook().GetMarket()[1].Len(),
		engine.GetOrderBook().GetHighestBid(),
		timeout,
	)
}
