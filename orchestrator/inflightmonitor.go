package orchestrator

import (
	"errors"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"log"
)

// InFlightMonitor keeps track of all work items which are started but not yet completed.
// It supports a single topic/partition pair.
// The completed offset is stored as a kafka offset in the main work consumer group.
// The started offset must be stored as an offset in a different consumer group.
//
// As work items are distributed to workers, the orchestrator ought to inform the in-flight monitor
// which work items have been distributed by calling the Start(offset) method. It is assumed the work
// started in the same order the kafka messages are stored in the partition.
//
// Workers need to inform the orchestrator when they complete the work by calling the Complete(offset)
// method. Because some work items might take more processing time, completion might happen in a completely
// different order than the one in which they were started.
//
// Each time work is completed, the offsets are stored in an array which is ordered by the kafka message offset.
// As soon as there is a common sequence of work which has been started and completed, the work consumer offset will
// be advanced by storing the offset in kafka.
type InFlightMonitor struct {
	started   *OffsetMonitor // Keeps track of items which were started, but no yet completed
	completed *OffsetMonitor // Keeps track of items which were completed, but not yet committed
}

func NewInFlightMonitor(workConsumer *SinglePartitionConsumer, startedOffsetConsumer *SinglePartitionConsumer) (*InFlightMonitor, error) {
	if startedOffsetConsumer.GroupID == workConsumer.GroupID {
		return nil, errors.New("expected consumers to use different consumer groups")
	}

	return &InFlightMonitor{
		completed: NewOffsetMonitor(workConsumer),
		started:   NewOffsetMonitor(startedOffsetConsumer),
	}, nil

}

// Start ought to be called every time an event is delivered to the worker.
func (i *InFlightMonitor) Start(offset kafka.TopicPartition) error {
	log.Printf("started work %v", offset)
	i.started.AddOffsets([]kafka.TopicPartition{offset})
	return i.advanceStartedOffset(offset)
}

// AckStartedOffsets ought to be called every time the started offset consumer receive  `kafka.OffsetsCommitted` event
func (i *InFlightMonitor) AckStartedOffset(offset kafka.TopicPartition) {
	log.Printf("successfully advanced started offset: %v", offset)
}

// Complete ought to be called every time a work is committed
func (i *InFlightMonitor) Complete(offset kafka.TopicPartition) error {
	// It is possible to receive a completion event from work managed by a previous execution of the work scheduler
	// which was not gracefully shutdown. Just ignore the event, because it will be reattempted.
	if !i.started.Contains(offset) {
		return nil
	}

	// Already expired or completed
	if i.completed.Contains(offset) {
		return nil
	}

	i.completed.AddOffsets([]kafka.TopicPartition{offset})
	log.Printf("completed work %v", offset)
	return i.commitFirstCompletedContinuousSequence()
}

func (i *InFlightMonitor) commitFirstCompletedContinuousSequence() error {
	log.Printf("all started: %v", i.started.GetOffsets().offsets)
	log.Printf("all completed: %v", i.completed.GetOffsets().offsets)
	common := i.completed.GetCommonStartingSequence(i.started)
	if len(common) == 0 {
		return nil
	}

	log.Printf("completed work sequence: %v", common)

	// just commit the last one
	maxOffset := common[len(common)-1]
	if err := i.completed.CommitOffset(maxOffset); err != nil {
		return err
	}
	log.Printf("successfully committed completed offsets %v", common)

	i.completed.RemoveCommittedOffsets(common)
	i.started.RemoveCommittedOffsets(common)
	log.Printf("removed completed items %v", common)

	return nil
}

func (i *InFlightMonitor) advanceStartedOffset(offset kafka.TopicPartition) error {
	log.Printf("advancing started offset to %v", offset)
	return i.started.CommitOffset(offset)
}

func (i *InFlightMonitor) HasMore() bool {
	return i.started.Count() > i.completed.Count()
}
