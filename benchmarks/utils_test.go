package benchmarks_test

import (
	"bufio"
	"context"
	"encoding/base64"
	"flag"
	"log"
	"math"
	"math/rand"
	"os"

	"github.com/segmentio/kafka-go"
	"gitlab.com/around25/products/matching-engine/model"
	"gitlab.com/around25/products/matching-engine/net"
)

const KAFKA_BROKER = "kafka:9092"
const KAFKA_CONSUMER_MARKET = "tstabc1"
const KAFKA_CONSUMER_TOPIC = "engine.orders." + KAFKA_CONSUMER_MARKET

var FLAG_GENERATE_DATA bool = false
var FLAG_GENERATE_KAFKA_DATA bool = false

func init() {
	flag.BoolVar(&FLAG_GENERATE_DATA, "gen", false, "Generate sample data")
	flag.BoolVar(&FLAG_GENERATE_KAFKA_DATA, "kafka", false, "Generate kafka sample data")
}

const BENCHMARK_TEST_FILE = "../priv/data/btcusd-2000000-limit-10-market.txt"

func GenerateOrdersInKafka(broker, topic string, n int) {
	GenerateRandomRecordsInFile(BENCHMARK_TEST_FILE, KAFKA_CONSUMER_MARKET, 2000000)
	if !FLAG_GENERATE_KAFKA_DATA {
		log.Println("No kafka sample data is going to be generated")
		return
	}

	producer := net.NewKafkaProducer([]string{broker}, topic)
	err := producer.Start()

	if err != nil {
		log.Println(err)
	}

	fh, err := os.Open(BENCHMARK_TEST_FILE)
	if err != nil {
		panic(err.Error())
	}
	defer fh.Close()
	bf := bufio.NewReader(fh)

	batch := make([]kafka.Message, 0, 20000)

	for j := 0; j < n; j++ {
		msg, _, err := bf.ReadLine()
		if err != nil {
			log.Fatalln(err)
		}
		data := make([]byte, 1000)
		base64.StdEncoding.Decode(data, msg)
		// err = producer.WriteMessages(context.Background(), kafka.Message{
		// 	Value: data,
		// })
		// if err != nil {
		// 	log.Fatal(err)
		// }
		batch = append(batch, kafka.Message{
			Value: data,
		})
		if len(batch) == cap(batch) {
			err = producer.WriteMessages(context.Background(), batch...)
			if err != nil {
				log.Fatal(err)
			}
			batch = make([]kafka.Message, 0, 20000)
		}
	}
	log.Println("Generated ", n, "orders in ", topic)
	producer.Close()
}

// GenerateRandomRecordsInFile create N orders and stores them in the given file
// Use "~/Incubator/go/src/engine/priv/data/market.txt" locally
// The file is opened with append and 0644 permissions
func GenerateRandomRecordsInFile(file string, market string, n int) bool {
	if !FLAG_GENERATE_DATA {
		log.Println("No local sample data will be generated in file: " + file)
		return false
	}

	rand.Seed(42)
	if _, err := os.Stat(file); err == nil {
		return false
	} else if !os.IsNotExist(err) {
		log.Fatal(err)
		return false
	}
	// create the file and generate the data
	fh, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		panic(err.Error())
	}
	defer fh.Close()
	for i := 0; i < n; i++ {
		id := uint64(i + 1)
		price := uint64(10000100-3*i-int(math.Ceil(10000*rand.Float64()))) * 100000000
		amount := uint64(10001-int(math.Ceil(10000*rand.Float64()))) * 100000000
		funds := price / 100000000 * amount
		side := model.MarketSide_Sell
		if int32(rand.Intn(2)%2) == 0 {
			side = model.MarketSide_Buy
		}
		evType := model.OrderType_Limit
		if int32(rand.Intn(100)%100) == 0 {
			evType = model.OrderType_Market
		}
		order := &model.Order{
			ID:        id,
			Market:    market,
			Amount:    amount,
			Price:     price,
			Side:      side,
			Type:      evType,
			EventType: model.CommandType_NewOrder,
			Stop:      model.StopLoss_None,
			StopPrice: 0,
			Funds:     funds,
		}
		data, _ := order.ToBinary()
		str := base64.StdEncoding.EncodeToString(data)
		fh.WriteString(str)
		if i < n {
			fh.WriteString("\n")
		}
	}
	return true
}
