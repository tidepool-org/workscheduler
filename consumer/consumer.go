package consumer

import (
	"context"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"log"
)

const (
	DefaultTimeoutMs = 100
)

type Consumer interface {
	Poll(context.Context) (*kafka.Message, error)
	StoreOffset(ctx context.Context, partition kafka.TopicPartition) error
}

type consumer struct {
	kafkaConsumer *kafka.Consumer
}

func New(kafkaConsumer *kafka.Consumer) (Consumer, error) {
	return &consumer{
		kafkaConsumer: kafkaConsumer,
	}, nil
}

func (c *consumer) Poll(ctx context.Context) (*kafka.Message, error) {
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			e := c.kafkaConsumer.Poll(DefaultTimeoutMs)
			if e == nil {
				continue
			}

			switch event := e.(type) {
			case *kafka.Message:
				if event.TopicPartition.Error != nil {
					return event, event.TopicPartition.Error
				}
				return event, nil
			case *kafka.Error:
				if event.IsFatal() {
					log.Fatalf(event.Error())
				}
				// ignore other errors - the client will try recover automatically
			default:
				log.Printf("ignored kafka event %v", event)
			}
		}
	}
}

func (c *consumer) StoreOffset(ctx context.Context, partition kafka.TopicPartition) error {
	_, err := c.kafkaConsumer.StoreOffsets([]kafka.TopicPartition{partition})
	return err
}
