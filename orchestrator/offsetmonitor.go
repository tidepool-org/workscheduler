package orchestrator

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type OffsetMonitor struct {
	consumer *SinglePartitionConsumer
	offsets *MinOrderedOffsets
}

// NewOffsetMonitor is used to monitor which topic-partition offsets
// are currently being worked on
func NewOffsetMonitor(consumer *SinglePartitionConsumer) *OffsetMonitor {
	return &OffsetMonitor{
		consumer: consumer,
		offsets: NewOrderedOffsets(),
	}
}

func (o *OffsetMonitor) AddOffsets(offsets []kafka.TopicPartition) {
	for _, offset := range offsets {
		o.offsets.Insert(offset)
	}
}

func (o *OffsetMonitor) CommitOffset(offset kafka.TopicPartition) error {
	// we need to commit 'n+1'
	next := kafka.TopicPartition{
		Topic: offset.Topic,
		Partition: offset.Partition,
		Offset: offset.Offset + 1,
	}

	_, err := o.consumer.CommitOffsets([]kafka.TopicPartition{next})
	return err
}

func (o *OffsetMonitor) GetCommonStartingSequence(other *OffsetMonitor) []kafka.TopicPartition {
	return o.offsets.GetCommonStartingSequence(other.GetOffsets())
}

func (o *OffsetMonitor) GetOffsets() *MinOrderedOffsets {
	return o.offsets
}

func (o *OffsetMonitor) RemoveCommittedOffsets(committedOffsets []kafka.TopicPartition) []kafka.TopicPartition {
	removed := make([]kafka.TopicPartition, len(committedOffsets))
	for _, committedOffset := range committedOffsets {
		offset := o.offsets.Peek()
		if !offsetsAreEqual(committedOffset, offset) {
			panic(fmt.Sprintf("trying to remove an offset which is out of order or is not tracked. Got %v, expected %v", committedOffset, offset))
		}

		removed = append(removed, o.offsets.Pop())
	}

	return removed
}

func (o *OffsetMonitor) RemoveOffsetsSmallerThanOrEqualTo(offset kafka.TopicPartition) []kafka.TopicPartition {
	removed := make([]kafka.TopicPartition, 0)
	for o.Count() > 0 {
		el := o.offsets.Peek()
		if !topicAndPartitionAreEqual(el, offset) {
			panic("unexpected topic and partition")
		}
		if el.Offset > offset.Offset {
			break
		}
		removed = append(removed, o.offsets.Pop())
	}

	return removed
}

func (o *OffsetMonitor) Contains(offset kafka.TopicPartition) bool {
	return o.offsets.Contains(offset)
}

func (o *OffsetMonitor) Count() int {
	return o.offsets.Len()
}

func offsetsAreEqual(a kafka.TopicPartition, b kafka.TopicPartition) bool {
	return a.Offset == b.Offset && topicAndPartitionAreEqual(a, b)
}

func topicAndPartitionAreEqual(a kafka.TopicPartition, b kafka.TopicPartition) bool {
	return a.Partition == b.Partition &&
		*a.Topic == *b.Topic
}
