package kafkaClient_test

import (
	"context"
	"github.com/qiaogy91/ioc/default/kafkaClient"
	"github.com/segmentio/kafka-go"
	"testing"
)

var (
	c = &kafkaClient.Client{
		Username: "adminscram",
		Password: "admin-secret",
		Debug:    true,
		Async:    false,
		Offset:   -2,
		Brokers:  []string{"127.0.0.1:9092"},
	}
	ctx = context.Background()
)

func TestClient_Producer(t *testing.T) {
	w := c.Producer("maudit")
	defer func() {
		_ = w.Close()
	}()

	if err := w.WriteMessages(ctx, kafka.Message{Value: []byte("hello world")}); err != nil {
		t.Fatal(err)
	}
}
func TestClient_Consumer(t *testing.T) {
	r := c.Consumer("topic01", "group01")

	for {
		msg, err := r.ReadMessage(ctx)
		if err != nil {
			t.Fatal(err)
			return
		}
		t.Logf("recived: %s", msg.Value)
	}
}
