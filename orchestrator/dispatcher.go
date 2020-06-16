package orchestrator

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/tidepool-org/workscheduler/workscheduler"
	"log"
	"sync"
)

type WorkDispatcher struct {
	consumer             *SinglePartitionConsumer
	workChan             chan workscheduler.Work
}

func NewWorkDispatcher(consumer *SinglePartitionConsumer) (*WorkDispatcher, error) {
	return &WorkDispatcher{
		consumer:             consumer,
		workChan:             make(chan workscheduler.Work),
	}, nil
}

func (w *WorkDispatcher) WorkChannel() <-chan workscheduler.Work {
	return w.workChan
}

func (w *WorkDispatcher) Run(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

eventLoop:
	for {
		select {
		// We can receive a shutdown signal while we're waiting for kafka messages
		case <-ctx.Done():
			fmt.Printf("Work disaptcher context has been terminated.")
			break eventLoop
		case ev := <-w.consumer.Events():
			switch event := ev.(type) {
			case *kafka.Message:
				work, err := w.msgToWork(event)
				if err != nil {
					log.Printf("Error while converting kafka message to work: %v", err)
					break
				}
				select {
				// We can also receive a shutdown signal while we're blocked on the work channel
				case <-ctx.Done():
					break eventLoop
				case w.workChan <- work:
					// Do nothing
				}
			case *kafka.Error:
				log.Printf("Received fatal error. Shutting down ...: %v", event.Error())
				break eventLoop
			}
		}
	}

	w.shutdown()
}

func (w *WorkDispatcher) msgToWork(msg *kafka.Message) (work workscheduler.Work, err error) {
	workSource, err := json.Marshal(msg.TopicPartition)
	if err != nil {
		return
	}

	work = workscheduler.Work{
		Source: &workscheduler.WorkSource{
			Source: string(workSource),
		},
		Data: msg.Value,
	}
	return
}

func (w *WorkDispatcher) shutdown() {
	close(w.workChan)
}
