package orchestrator

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"sort"
	"time"
)

type expirationRecord struct {
	offset         kafka.TopicPartition
	expirationTime time.Time
}

type ExpirationMonitor struct {
	records []*expirationRecord // kafka message offsets ordered by expiration time
	timeout time.Duration
}

func NewExpirationMonitor(timeout time.Duration) *ExpirationMonitor {
	return &ExpirationMonitor{
		records: make([]*expirationRecord, 0),
		timeout: timeout,
	}
}

func (e *ExpirationMonitor) Add(offset kafka.TopicPartition) {
	now := time.Now()
	expirationTime := now.Add(e.timeout)
	e.records = append(e.records, &expirationRecord{
		offset:         offset,
		expirationTime: expirationTime,
	})
}

func (e *ExpirationMonitor) Remove(offset kafka.TopicPartition) {
	// TODO: this assumes offsets are always increasing, which is the case atm, but this might change.
	// The array is actually ordered by the expiration time, but it happens to be ordered by offsets too,
	// because the offsets are added for expiration monitoring in the order they are started, which is
	// always the order they are stored in the kafka partition (i.e. with monotonically increasing offsets).
	i := sort.Search(len(e.records), func(i int) bool { return e.records[i].offset.Offset >= offset.Offset })
	if i < len(e.records) && offsetsAreEqual(e.records[i].offset, offset) {
		e.records = append(e.records[:i], e.records[i+1:]...)
	}
}

func (e *ExpirationMonitor) RemoveExpired() []kafka.TopicPartition {
	expired := make([]kafka.TopicPartition, 0)
	firstNonExpired := len(e.records)
	now := time.Now()
	for i, rec := range e.records {
		if rec.expirationTime.After(now) {
			firstNonExpired = i
			break
		}

		expired = append(expired, rec.offset)
	}

	e.records = e.records[firstNonExpired:]

	return expired
}
