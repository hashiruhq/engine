package trading_engine_test

import (
	"fmt"
	"log"
	"math/rand"
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

	for i := 0; i < n; i++ {
		rnd := rand.Float64()
		msg := fmt.Sprintf("order[%d][ %s ] %d - %f @ %f", 1, "id", 1+rand.Intn(2)%2, 10001.0-10000*rand.Float64(), 4000100.00-(float64)(i)-1000000*rnd)
		_, _, err := producer.Send(([]byte)(msg))
		if err != nil {
			log.Println(err)
			break
		}
	}
	producer.Close()
}
