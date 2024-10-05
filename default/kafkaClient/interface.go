package kafkaClient

import (
	"github.com/qiaogy91/ioc"
	"github.com/segmentio/kafka-go"
)

const AppName = "kafkaClient"

func GetConsumer(topic, groupID string) *kafka.Reader {
	return ioc.Default().Get(AppName).(*Client).Consumer(topic, groupID)
}

func GetProducer(topic string) *kafka.Writer {
	return ioc.Default().Get(AppName).(*Client).Producer(topic)
}
