package benchmarks_test

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/golang/protobuf/proto"
	"gitlab.com/around25/products/matching-engine/engine"
)

var arr []engine.Order = make([]engine.Order, 0, 2000000)
var msgs [][]byte = make([][]byte, 0, 2000000)
var ngin = engine.NewTradingEngine()

func init() {
	rand.Seed(42)
	testFile := "/Users/cosmin/Incubator/go/src/gitlab.com/around25/products/matching-engine/priv/data/market.txt"
	// GenerateRandomRecordsInFile(&testFile, 2000000)
	fh, err := os.Open(testFile)
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
		order := engine.Order{}
		order.FromBinary(data)
		data, _ = order.ToBinary()
		msgs = append(msgs, data)
		arr = append(arr, order)
	}
}

func BenchmarkDecodeFromProto(benchmark *testing.B) {
	order := &engine.Order{
		ID:        "TST_1",
		Market:    "btc-usd",
		Amount:    848382829993942,
		Price:     131221300010201,
		Side:      engine.MarketSide_Buy,
		Type:      engine.OrderType_Limit,
		EventType: engine.CommandType_NewOrder,
		Stop:      engine.StopLoss_Loss,
		StopPrice: 1313231100010201,
		Funds:     10100010133232313,
	}
	binary, _ := proto.Marshal(order)
	for i := 0; i < benchmark.N; i++ {
		proto.Unmarshal(binary, order)
	}
}

func BenchmarkEncodeToProto(benchmark *testing.B) {
	order := &engine.Order{
		ID:        "TST_1",
		Market:    "btc-usd",
		Amount:    848382829993942,
		Price:     131221300010201,
		Side:      engine.MarketSide_Buy,
		Type:      engine.OrderType_Limit,
		EventType: engine.CommandType_NewOrder,
		Stop:      engine.StopLoss_Loss,
		StopPrice: 1313231100010201,
		Funds:     10100010133232313,
	}
	for i := 0; i < benchmark.N; i++ {
		proto.Marshal(order)
	}
}

func BenchmarkWithRandomData(benchmark *testing.B) {
	startTime := time.Now().UnixNano()
	// ngin := engine.NewTradingEngine()
	for j := 0; j < benchmark.N; j++ {
		ngin.Process(arr[j])
	}
	PrintOrderLogs(ngin, benchmark.N, startTime)
}

func BenchmarkWithDecodeAndEncodeRandomData(benchmark *testing.B) {
	startTime := time.Now().UnixNano()
	// ngin := engine.NewTradingEngine()
	for j := 0; j < benchmark.N; j++ {
		order := engine.Order{}
		order.FromBinary(msgs[j])
		trades := ngin.Process(order)
		for _, trade := range trades {
			trade.ToBinary()
		}
	}
	PrintOrderLogs(ngin, benchmark.N, startTime)
}

// GenerateRandomRecordsInFile create N orders and stores them in the given file
// Use "~/Incubator/go/src/engine/priv/data/market.txt" locally
// The file is opened with append and 0644 permissions
func GenerateRandomRecordsInFile(file *string, n int) {
	rand.Seed(42)
	fh, err := os.OpenFile(*file, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		panic(err.Error())
	}
	defer fh.Close()
	for i := 0; i < n; i++ {
		id := "ID_" + strconv.Itoa(i)
		price := uint64(10000100-3*i-int(math.Ceil(10000*rand.Float64()))) * 100000000
		amount := uint64(10001-int(math.Ceil(10000*rand.Float64()))) * 100000000
		side := engine.MarketSide_Sell
		if int32(rand.Intn(2)%2) == 0 {
			side = engine.MarketSide_Buy
		}
		order := &engine.Order{
			ID:        id,
			Market:    "btc-usd",
			Amount:    amount,
			Price:     price,
			Side:      side,
			Type:      engine.OrderType_Limit,
			EventType: engine.CommandType_NewOrder,
			Stop:      engine.StopLoss_Loss,
			StopPrice: 1313231100010201,
			Funds:     10100010133232313,
		}
		data, _ := order.ToBinary()
		str := base64.StdEncoding.EncodeToString(data)
		fh.WriteString(str)
		if i < n {
			fh.WriteString("\n")
		}
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
