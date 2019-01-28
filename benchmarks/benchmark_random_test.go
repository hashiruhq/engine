package benchmarks_test

import (
	"bufio"
	"log"
	"math/rand"
	"os"
	"testing"

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
	// startTime := time.Now().UnixNano()
	ngin := engine.NewTradingEngine()
	for j := 0; j < benchmark.N; j++ {
		ngin.Process(arr[j])
	}
	// utils.PrintOrderLogs(ngin, benchmark.N, startTime)
}
