package kafka

import (
	"context"
	"github.com/qiaogy91/ioc"
	"github.com/qiaogy91/ioc/config/log"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl"
	"github.com/segmentio/kafka-go/sasl/scram"
	"log/slog"
	"time"
)

/*
注意：当异步模式开启时，Producer 必须持续至少运行 BatchTimeout 时间后再结束，确保消息能够在 BatchTimeout 时间到达后写入 Kafka 集群
*/

type Client struct {
	ioc.ObjectImpl
	log          *slog.Logger
	Username     string   `json:"username" yaml:"username"`
	Password     string   `json:"password" yaml:"password"`
	Brokers      []string `json:"brokers" yaml:"brokers"`
	Async        bool     `json:"async" yaml:"async"`
	BatchTimeout int      `json:"batchTimeout" yaml:"batchTimeout"`
	Offset       int64    `json:"offset" yaml:"offset"`
	mechanism    sasl.Mechanism
}

func (c *Client) Name() string  { return AppName }
func (c *Client) Priority() int { return 203 }
func (c *Client) Init() {
	c.log = log.Sub(AppName)
	if c.Username == "" {
		return
	}

	mechanism, err := scram.Mechanism(scram.SHA512, c.Username, c.Password)
	if err != nil {
		panic(err)
	}
	c.mechanism = mechanism
}
func (c *Client) Close(ctx context.Context) error {
	c.log.Info("closed completed", slog.String("namespace", ioc.DefaultNamespace))
	return nil
}

func (c *Client) Producer(topic string) *kafka.Writer {
	return &kafka.Writer{
		Addr:                   kafka.TCP(c.Brokers...),
		Topic:                  topic,
		Balancer:               &kafka.LeastBytes{},
		Transport:              &kafka.Transport{SASL: c.mechanism},
		AllowAutoTopicCreation: true,
		Async:                  c.Async,
		BatchTimeout:           time.Duration(c.BatchTimeout) * time.Second,
	}
}

func (c *Client) Consumer(topic, groupId string) *kafka.Reader {
	conf := kafka.ReaderConfig{
		Brokers:     c.Brokers,
		Topic:       topic,
		GroupID:     groupId,
		MaxBytes:    10e6,
		StartOffset: c.Offset,
		Dialer:      &kafka.Dialer{Timeout: 10 * time.Second, SASLMechanism: c.mechanism},
	}
	return kafka.NewReader(conf)
}

func init() {
	ioc.Default().Registry(&Client{})
}
