package benchmarks_test

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"strconv"
	"testing"
	"time"

	"gitlab.com/around25/products/matching-engine/engine"
)

var arr []engine.Order = make([]engine.Order, 0, 2000000)

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
		order := engine.Order{}
		order.FromJSON([]byte(msg))
		arr = append(arr, order)
	}
}

func BenchmarkWithRandomData(benchmark *testing.B) {
	startTime := time.Now().UnixNano()
	ngin := engine.NewTradingEngine()
	for j := 0; j < benchmark.N; j++ {
		ngin.Process(arr[j])
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
		price := 10000100 - 3*i - int(math.Ceil(10000*rand.Float64()))
		amount := 10001 - int(math.Ceil(10000*rand.Float64()))
		side := int8(1 + rand.Intn(2)%2)
		str := fmt.Sprintf(`{"base":"eth","quote":"btc","id":"%s","price":"%d","amount":"%d","side":%d,"type":1,"event_type":1}`+"\n", id, price, amount, side)
		fh.WriteString(str)
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
