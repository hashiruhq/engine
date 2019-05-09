package engine_test

import (
	"context"
	"fmt"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/segmentio/kafka-go"
	"gitlab.com/around25/products/matching-engine/net"
)

func generateOrdersInKafka(n int) {
	rand.Seed(42)
	kafkaBroker := "kafka:9092"
	kafkaOrderTopic := "trading.order.btc.eth"

	producer := net.NewKafkaProducer([]string{kafkaBroker})
	err := producer.Start()

	go func(producer net.KafkaProducer) {
		errors := producer.Errors()
		for err := range errors {
			value, _ := err.Msg.Value.Encode()
			log.Print("Error received from trades producer ", (string)(value), err)
		}
	}(producer)
	if err != nil {
		log.Println(err)
	}

	for i := 0; i < n; i++ {
		id := "ID_" + fmt.Sprintf("%d", i)
		price := 4000100 - 3*i - int(math.Ceil(10000*rand.Float64()))
		amount := 10001 - int(math.Ceil(10000*rand.Float64()))
		side := int8(1 + rand.Intn(2)%2)
		producer.WriteMessages(context.Background(), kafka.Message{
			Value: ([]byte)(fmt.Sprintf(`{"base":"eth","quote":"btc","id":"%s","price":"%d","amount":"%d","side":%d,"type":1}`, id, price, amount, side)),
		})
	}

	time.Sleep(time.Millisecond * 300)

	producer.Close()
}
