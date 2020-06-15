package orchestrator

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/pkg/errors"
)

type SinglePartitionConsumer struct {
	*kafka.Consumer
	tp kafka.TopicPartition

	GroupID string
}

func NewSinglePartitionConsumer(brokers, topic, groupId string) (*SinglePartitionConsumer, error) {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":        brokers,
		"group.id":                 groupId,
		"go.events.channel.enable": true,
		"enable.auto.offset.store": false, // commit offsets manually
		"enable.auto.commit":       false,
	})
	if err != nil {
		return nil, err
	}

	meta, err := consumer.GetMetadata(&topic, true, 5000)
	if err != nil {
		return nil, err
	}

	topicMetadata, ok := meta.Topics[topic]
	if !ok {
		return nil, errors.New("topic doesn't exist")
	}

	if count := len(topicMetadata.Partitions); count != 1 {
		return nil, errors.New(fmt.Sprintf("expected a topic with a single partition, got topic with %v partitions", count))
	}

	if err = consumer.Subscribe(topic, nil); err != nil {
		return nil, errors.Wrap(err, "could not subscribe to topic")
	}

	tp := kafka.TopicPartition{
		Topic:     &topic,
		Partition: topicMetadata.Partitions[0].ID,
	}

	return &SinglePartitionConsumer{
		Consumer: consumer,
		GroupID:  groupId,
		tp:       tp,
	}, nil
}

func (s *SinglePartitionConsumer) Pause() error {
	return s.Consumer.Pause([]kafka.TopicPartition{s.tp})
}
