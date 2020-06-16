package orchestrator

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/pkg/errors"
	"github.com/tidepool-org/workscheduler/workscheduler"
	"log"
	"sync"
	"time"
)

type Orchestrator interface {
	WorkChannel() <-chan workscheduler.Work              // channel to read work requests from a work provider
	WorkStartedChannel() chan<- workscheduler.WorkSource // channel to write when work is started
	ResponseChannel() chan<- workscheduler.WorkSource    // channel to write work responses to a work provider
	Run(context.Context, *sync.WaitGroup)
}

type kafkaSinglePartitionWorkOrchestrator struct {
	workDispatcher    *WorkDispatcher
	workConsumer      *SinglePartitionConsumer
	inFlightMonitor   *InFlightMonitor
	expirationMonitor *ExpirationMonitor
	responseChan      chan workscheduler.WorkSource
	workStartedChan   chan workscheduler.WorkSource
}

type KafkaWorkOrchestratorParams struct {
	Brokers     string
	Prefix      string
	Topic       string
	WorkTimeout time.Duration
}

func (k KafkaWorkOrchestratorParams) GetPrefixedTopic() string {
	return fmt.Sprintf("%s%s", k.Prefix, k.Topic)
}

func NewKafkaWorkOrchestrator(params KafkaWorkOrchestratorParams) (Orchestrator, error) {
	prefixedTopic := params.GetPrefixedTopic()
	workConsumer, err := NewSinglePartitionConsumer(params.Brokers, prefixedTopic, "work_consumer")
	if err != nil {
		return nil, errors.Wrap(err, "failed to start work consumer")
	}
	startedOffsetConsumer, err := NewSinglePartitionConsumer(params.Brokers, prefixedTopic, "_started_offset_tracker_new")
	if err != nil {
		return nil, errors.Wrap(err, "failed to start _started_offset_tracker")
	}
	inFlightMonitor, err := NewInFlightMonitor(workConsumer, startedOffsetConsumer)
	if err != nil {
		return nil, err
	}

	workDispatcher, err := NewWorkDispatcher(workConsumer)
	if err != nil {
		return nil, errors.Wrap(err, "failed to start work dispatcher")
	}

	return &kafkaSinglePartitionWorkOrchestrator{
		workConsumer:      workConsumer,
		workDispatcher:    workDispatcher,
		inFlightMonitor:   inFlightMonitor,
		expirationMonitor: NewExpirationMonitor(params.WorkTimeout),
		responseChan:      make(chan workscheduler.WorkSource),
		workStartedChan:   make(chan workscheduler.WorkSource),
	}, nil
}

func (k *kafkaSinglePartitionWorkOrchestrator) WorkChannel() <-chan workscheduler.Work {
	return k.workDispatcher.WorkChannel()
}

func (k *kafkaSinglePartitionWorkOrchestrator) WorkStartedChannel() chan<- workscheduler.WorkSource {
	return k.workStartedChan
}

func (k *kafkaSinglePartitionWorkOrchestrator) ResponseChannel() chan<- workscheduler.WorkSource {
	return k.responseChan
}

func (k *kafkaSinglePartitionWorkOrchestrator) Run(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	ticker := time.NewTicker(time.Second)

	wg.Add(1)
	wdCtx, _ := context.WithCancel(ctx)
	go k.workDispatcher.Run(wdCtx, wg)

	for {
		select {
		case <-ctx.Done():
			return
		case workSource := <-k.workStartedChan:
			k.start(workSource)
		case workSource := <-k.responseChan:
			k.complete(workSource)
		case <-ticker.C:
			// TODO: move the ticker to the expiration monitor
			k.expire()
		}
	}
}

func (k *kafkaSinglePartitionWorkOrchestrator) start(source workscheduler.WorkSource) {
	offset, err := k.getOffset(source)
	if err != nil {
		log.Printf("unable to get kafka offset from work source %v: %v", source, err.Error())
		return
	}

	k.expirationMonitor.Add(offset)
	if err := k.inFlightMonitor.Start(offset); err != nil {
		log.Printf("error while marking work as started: %v", err)
	}
}

func (k *kafkaSinglePartitionWorkOrchestrator) complete(src workscheduler.WorkSource) {
	offset, err := k.getOffset(src)
	if err != nil {
		log.Printf("could not get offset from work source %v", src)
		return
	}

	k.expirationMonitor.Remove(offset)
	if err = k.inFlightMonitor.Complete(offset); err != nil {
		log.Printf("could not complete work in in-flight monitor %v", err)
	}
}

func (k *kafkaSinglePartitionWorkOrchestrator) expire() {
	expired := k.expirationMonitor.RemoveExpired()
	if len(expired) == 0 {
		return
	}

	log.Printf("offets expired: %v", expired)
	for _, offset := range expired {
		if err := k.inFlightMonitor.Complete(offset); err != nil {
			log.Printf("could not expire offset: %v", err)
		}
	}
}

func (k kafkaSinglePartitionWorkOrchestrator) getOffset(src workscheduler.WorkSource) (offset kafka.TopicPartition, err error) {
	err = json.Unmarshal([]byte(src.Source), &offset)
	return
}
