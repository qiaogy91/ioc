package kafkaClient

import (
	"context"
	"github.com/qiaogy91/ioc"
	"github.com/qiaogy91/ioc/config/log"
	"github.com/rs/zerolog"
	kafka "github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl"
	"github.com/segmentio/kafka-go/sasl/scram"
	"time"
)

type Client struct {
	ioc.ObjectImpl
	log       *zerolog.Logger
	Username  string   `json:"username" yaml:"username"`
	Password  string   `json:"password" yaml:"password"`
	Debug     bool     `json:"debug" yaml:"debug"`
	Brokers   []string `json:"brokers" yaml:"brokers"`
	Async     bool     `json:"async" yaml:"async"`
	Offset    int64    `json:"offset" yaml:"offset"`
	producers []*kafka.Writer
	consumers []*kafka.Reader
	mechanism sasl.Mechanism
}

func (c *Client) Name() string  { return AppName }
func (c *Client) Priority() int { return 201 }
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

func (c *Client) Producer(topic string) *kafka.Writer {
	w := &kafka.Writer{
		Addr:                   kafka.TCP(c.Brokers...),
		Topic:                  topic,
		Balancer:               &kafka.LeastBytes{},
		Transport:              &kafka.Transport{SASL: c.mechanism},
		AllowAutoTopicCreation: true,
		Async:                  c.Async,
		Completion: func(messages []kafka.Message, err error) {
			if c.Debug && err != nil {
				c.log.Error().Msgf("Producer write failed, %s", err)
			}
		},
	}
	c.producers = append(c.producers, w)
	return w
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

	r := kafka.NewReader(conf)
	c.consumers = append(c.consumers, r)
	return r
}

func (c *Client) Close(ctx context.Context) error {
	for _, item := range c.producers {
		if err := item.Close(); err != nil {
			c.log.Error().Msgf("close producer failed, %s", err)
		}
	}

	for _, item := range c.consumers {
		if err := item.Close(); err != nil {
			c.log.Error().Msgf("close consumer failed, %s", err)
		}
	}

	return nil
}
func init() {
	ioc.Default().Registry(&Client{})
}
