package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/segmentio/kafka-go"
	"log/slog"
	"time"
)

// Client - клиент продьюсера для Kafka
type Client struct {
	Writer *kafka.Writer
	log    *slog.Logger
}

type Message struct {
	Data   []byte    `json:"data"`
	SendOn time.Time `json:"send_on"`
}

// New создает и инициализирует клиента продюьсера для Kafka.
func New(broker string, topic string, log *slog.Logger) (*Client, error) {
	if broker == "" || topic == "" {
		return nil, errors.New("не указаны параметры подключения к Kafka")
	}

	c := Client{
		log: log,
	}

	c.Writer = &kafka.Writer{
		Addr:                   kafka.TCP(broker),
		Topic:                  topic,
		Balancer:               &kafka.LeastBytes{},
		AllowAutoTopicCreation: true,
	}

	return &c, nil
}

func (c *Client) Produce(msgVal []byte, ctx context.Context) {
	m := Message{
		Data:   msgVal,
		SendOn: time.Now(),
	}

	marshal, err := json.Marshal(m)
	if err != nil {
		c.log.Error(err.Error())
		return
	}

	err = c.Writer.WriteMessages(ctx, kafka.Message{
		Value: marshal,
	})
	if err != nil {
		c.log.Error(err.Error())
		return
	}

	c.log.Debug("Produce " + string(msgVal))
}
