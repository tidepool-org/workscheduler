package orchestrator

import (
	"context"
	"errors"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"log"
	"sync"
	"time"
)

type WorkStatus int

const (
	WorkStarted WorkStatus = iota
	WorkStartedAndCommitted
	WorkCompleted
	WorkCompletedAndCommitted
)

type WorkUpdate struct {
	Offset kafka.TopicPartition
	Status WorkStatus
}

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

	lastStartedAndCommitted *kafka.TopicPartition
	workUpdatesChan         chan *WorkUpdate
}

func NewInFlightMonitor(workConsumer *SinglePartitionConsumer, startedOffsetConsumer *SinglePartitionConsumer) (*InFlightMonitor, error) {
	if startedOffsetConsumer.GroupID == workConsumer.GroupID {
		return nil, errors.New("expected consumers to use different consumer groups")
	}

	workUpdatesChan := make(chan *WorkUpdate)

	return &InFlightMonitor{
		completed:       NewOffsetMonitor(workConsumer),
		started:         NewOffsetMonitor(startedOffsetConsumer),
		workUpdatesChan: workUpdatesChan,
	}, nil
}

func (i *InFlightMonitor) GetWorkUpdatesChannel() chan<- *WorkUpdate {
	return i.workUpdatesChan
}

func (i *InFlightMonitor) Run(ctx context.Context, wg *sync.WaitGroup) {
	t := time.NewTicker(time.Second * 30)

	for {
		select {
		case <-ctx.Done():
			fmt.Printf("done")
		case update := <-i.workUpdatesChan:
			i.handleWorkUpdate(update)
		case <-t.C:
			i.commit()
		}
	}
}

func (i *InFlightMonitor) handleWorkUpdate(update *WorkUpdate) {
	switch update.Status {
	case WorkStarted:
		i.start(update.Offset)
	case WorkCompleted:
		i.complete(update.Offset)
	case WorkStartedAndCommitted:
		i.handleStartedWorkCommit(update.Offset)
	case WorkCompletedAndCommitted:
		i.handleCompletedWorkCommit(update.Offset)
	}
}

func (i *InFlightMonitor) handleStartedWorkCommit(offset kafka.TopicPartition) {
	// commits might be completed out of order
	if i.lastStartedAndCommitted == nil || i.lastStartedAndCommitted.Offset < offset.Offset {
		i.lastStartedAndCommitted = &offset
	}
}

func (i *InFlightMonitor) commit() {
	if i.shouldCommitStarted() {
		i.commitStarted()
	}
	if i.shouldCommitCompleted() {
		i.commitFirstCompletedContinuousSequence()
	}
}

func (i *InFlightMonitor) commitStarted() {
	go func() {
		offset := i.started.Last()
		if err := i.started.CommitOffset(*offset); err != nil {
			log.Println(fmt.Sprintf("error committing started work offset: %v", err))
			return
		}

		i.workUpdatesChan <- &WorkUpdate{
			Offset: *offset,
			Status: WorkStartedAndCommitted,
		}
	}()
}

func (i *InFlightMonitor) shouldCommitStarted() bool {
	return i.started.Last() != i.lastStartedAndCommitted
}

func (i *InFlightMonitor) shouldCommitCompleted() bool {
	return len(i.completed.GetCommonStartingSequence(i.started)) == 0
}

func (i *InFlightMonitor) start(offset kafka.TopicPartition) {
	log.Printf("started work %v", offset)
	i.started.AddOffsets([]kafka.TopicPartition{offset})
}

// Complete ought to be called every time a work is committed
func (i *InFlightMonitor) complete(offset kafka.TopicPartition) {
	// It is possible to receive a completion event from work managed by a previous execution of the work scheduler
	// which was not gracefully shutdown. Just ignore the event, because it will be reattempted.
	if !i.started.Contains(offset) {
		return
	}

	// Already expired or completed
	if i.completed.Contains(offset) {
		return
	}

	i.completed.AddOffsets([]kafka.TopicPartition{offset})
	log.Printf("completed work %v", offset)
}

func (i *InFlightMonitor) commitFirstCompletedContinuousSequence() {
	log.Printf("all started: %v", i.started.GetOffsets().offsets)
	log.Printf("all completed: %v", i.completed.GetOffsets().offsets)
	common := i.completed.GetCommonStartingSequence(i.started)
	if len(common) == 0 {
		return
	}

	log.Printf("completed work sequence: %v", common)

	// just commit the last one
	maxOffset := common[len(common)-1]
	go func() {
		if err := i.completed.CommitOffset(maxOffset); err != nil {
			log.Println(fmt.Sprintf("error committing completed work offset: %v", err))
			return
		}

		i.workUpdatesChan <- &WorkUpdate{
			Offset: maxOffset,
			Status: WorkCompletedAndCommitted,
		}
	}()
}

func (i *InFlightMonitor) handleCompletedWorkCommit(offset kafka.TopicPartition) {
	log.Printf("successfully committed completed offset %v", offset)
	removedCompleted := i.completed.RemoveOffsetsSmallerThanOrEqualTo(offset)
	log.Printf("removed completed items %v", removedCompleted)
	removedStarted := i.started.RemoveOffsetsSmallerThanOrEqualTo(offset)
	log.Printf("removed started items %v", removedStarted)
}

func (i *InFlightMonitor) HasMore() bool {
	return i.started.Count() > i.completed.Count()
}
