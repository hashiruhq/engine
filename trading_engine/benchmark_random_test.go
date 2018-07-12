package trading_engine_test

import (
	"bufio"
	"log"
	"math/rand"
	"os"
	"testing"
	"trading_engine/trading_engine"
)

var arr []trading_engine.Order = make([]trading_engine.Order, 0, 1000000)

func init() {
	rand.Seed(42)
	testFile := "/Users/cosmin/Incubator/go/src/trading_engine/priv/data/market.txt"
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
		order := trading_engine.Order{}
		order.FromJSON([]byte(msg))
		arr = append(arr, order)
	}
}

func BenchmarkWithRandomData(benchmark *testing.B) {
	// startTime := time.Now().UnixNano()
	engine := trading_engine.NewTradingEngine()
	for j := 0; j < benchmark.N; j++ {
		engine.Process(arr[j])
	}
	// utils.PrintOrderLogs(engine, benchmark.N, startTime)
}
