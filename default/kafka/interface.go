package kafka

import (
	"github.com/qiaogy91/ioc"
	"github.com/segmentio/kafka-go"
)

const AppName = "kafka"

func GetClient() Service {
	return ioc.Default().Get(AppName).(Service)
}

type Service interface {
	Producer(topic string) *kafka.Writer
	Consumer(topic, groupId string) *kafka.Reader
}
