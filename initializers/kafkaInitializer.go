package initializers

import (
	"log"
	"os"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

var KafkaProducer *kafka.Producer
var ProductClickTopic string

func KafkaInitializer() {

	var err error
	KafkaProducer, err = kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost"})
	if err != nil {
		panic(err)
	}

	go messageHandler()
	ProductClickTopic = os.Getenv("PRODUCT_CLICK_TOPIC_NAME")

}

func messageHandler() {
	for e := range KafkaProducer.Events() {
		switch ev := e.(type) {
		case *kafka.Message:
			if ev.TopicPartition.Error != nil {
				log.Printf("Delivery failed: %v\n", ev.TopicPartition)
			} else {
				log.Printf("Delivered message to %v\n", ev.TopicPartition)
			}
		}
	}
}
