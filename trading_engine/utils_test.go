package trading_engine_test

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"strconv"
	"trading_engine/net"
)

func generateOrdersInKafka(n int) {
	rand.Seed(42)
	kafkaBroker := "kafka:9092"
	kafkaOrderTopic := "trading.order.btc.eth"

	producer := net.NewKafkaProducer([]string{kafkaBroker}, kafkaOrderTopic)
	err := producer.Start()
	if err != nil {
		log.Println(err)
	}
	defer producer.Close()

	genRandMsgs := func(base, n int) *[][]byte {
		msgs := make([][]byte, 0, n)
		for j := 0; j < n; j++ {
			id := "ID_" + fmt.Sprintf("%d", rand.Uint32())
			price := 4000100 - 3*(base+j) - int(math.Ceil(10000*rand.Float64()))
			amount := 10001 - int(math.Ceil(10000*rand.Float64()))
			side := int8(1 + rand.Intn(2)%2)
			msg := fmt.Sprintf(`{"base": "sym", "market": "tst", "id":"%s", "price": %d, "amount": %d, "side": %d, "category": 1}`, id, price, amount, side)
			msgs = append(msgs, ([]byte)(msg))
		}
		return &msgs
	}

	chunkSize := 1000
	chunks := n / chunkSize
	lastChunk := n % chunkSize
	for i := 0; i < chunks; i++ {
		msgs := genRandMsgs(i*chunkSize, chunkSize)
		err := producer.SendMessages(*msgs)
		if err != nil {
			log.Println(err)
			break
		}
	}
	msgs := genRandMsgs(chunks*chunkSize, lastChunk)
	err = producer.SendMessages(*msgs)
	if err != nil {
		log.Println(err)
	}
}

// GenerateRandomRecordsInFile create N orders and stores them in the given file
// Use "~/Incubator/go/src/trading_engine/priv/data/market.txt" locally
// The file is opened with append and 0644 permissions
func GenerateRandomRecordsInFile(file *string, n int) {
	rand.Seed(42)
	fh, err := os.OpenFile(*file, os.O_WRONLY, 0644)
	if err != nil {
		panic(err.Error())
	}
	defer fh.Close()
	for i := 0; i < n; i++ {
		id := "ID_" + strconv.Itoa(i)
		price := 10000100 - 3*i - int(math.Ceil(10000*rand.Float64()))
		amount := 10001 - int(math.Ceil(10000*rand.Float64()))
		side := int8(1 + rand.Intn(2)%2)
		str := fmt.Sprintf(`{"base": "sym", "market": "tst", "id":"%s", "price": %d, "amount": %d, "side": %d, "category": 1}`+"\n", id, price, amount, side)
		fh.WriteString(str)
	}
}
