package orchestrator

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"sort"
)

type MinOrderedOffsets struct {
	offsets []kafka.TopicPartition
}

func NewOrderedOffsets() *MinOrderedOffsets {
	return &MinOrderedOffsets{
		offsets: make([]kafka.TopicPartition, 0),
	}
}

func (o *MinOrderedOffsets) Insert(offset kafka.TopicPartition) {
	index := sort.Search(len(o.offsets), func(i int) bool { return o.offsets[i].Offset > offset.Offset })
	o.offsets = append(o.offsets, offset)
	copy(o.offsets[index+1:], o.offsets[index:])
	o.offsets[index] = offset
}

func (o *MinOrderedOffsets) Pop() kafka.TopicPartition {
	x := o.offsets[0]
	o.offsets = o.offsets[1:]
	return x
}

func (o *MinOrderedOffsets) Get(i int) kafka.TopicPartition {
	return o.offsets[i]
}

func (o *MinOrderedOffsets) Peek() kafka.TopicPartition {
	return o.offsets[0]
}

func (o *MinOrderedOffsets) Len() int {
	return len(o.offsets)
}

func (o *MinOrderedOffsets) Contains(offset kafka.TopicPartition) bool {
	// TODO: use binary search
	for _, tp := range o.offsets {
		if offsetsAreEqual(offset, tp) {
			return true
		}
	}

	return false
}

func (o *MinOrderedOffsets) GetCommonStartingSequence(other *MinOrderedOffsets) []kafka.TopicPartition {
	common := make([]kafka.TopicPartition, 0)
	for i, tp := range o.offsets {
		if i >= other.Len() || tp.Offset != other.Get(i).Offset {
			break
		}
		common = append(common, tp)
	}
	return common
}
