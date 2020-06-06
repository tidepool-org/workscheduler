package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/tidepool-org/workscheduler/consumer"
	pb "github.com/tidepool-org/workscheduler/workscheduler"
	"google.golang.org/grpc"
	"log"
	"net"
	"sync"
)

// Config is the configuration of the work scheduler
type Config struct {
	Brokers string
	Group   string
	Prefix  string
	Topic   string
	Port    string
}

// WorkSchedulerServer is used to implement workscheduler
type WorkSchedulerServer struct {
	pb.UnimplementedWorkSchedulerServer
	Config     Config
	grpcServer *grpc.Server
	consumer   consumer.Consumer
}

// NewWorkSchedulerServer create a new WorkSchedulerServer
func NewWorkSchedulerServer(c Config) *WorkSchedulerServer {
	grpcServer := grpc.NewServer()
	pb.RegisterWorkSchedulerServer(grpcServer, &WorkSchedulerServer{})

	return &WorkSchedulerServer{
		Config:     c,
		grpcServer: grpcServer,
	}
}

// Poll for work to process
func (s *WorkSchedulerServer) Poll(ctx context.Context, in *empty.Empty) (*pb.Work, error) {
	// blocks until there's work or until the context is terminated
	msg, err := s.consumer.Poll(ctx)
	if err != nil {
		return nil, err
	}

	ws, err := json.Marshal(msg.TopicPartition)
	if err != nil {
		return nil, err
	}

	work := &pb.Work{
		Source: &pb.WorkSource{
			Source: string(ws),
		},
		Data:       msg.Value,
	}

	return work, nil
}

// Ping to check for health of work scheduler server
func (s *WorkSchedulerServer) Ping(ctx context.Context, in *empty.Empty) (*empty.Empty, error) {
	return &empty.Empty{}, nil
}

// Failed to report failed job
func (s *WorkSchedulerServer) Failed(ctx context.Context, in *pb.WorkSource) (*empty.Empty, error) {
	return s.acknowledgeWork(ctx, in)
}

// Complete to report completed job
func (s *WorkSchedulerServer) Complete(ctx context.Context, in *pb.WorkSource) (*pb.WorkOutput, error) {
	_, err := s.acknowledgeWork(ctx, in)
	return &pb.WorkOutput{}, err
}

func (s *WorkSchedulerServer) acknowledgeWork(ctx context.Context, in *pb.WorkSource) (*empty.Empty, error) {
	tp := kafka.TopicPartition{}
	if err := json.Unmarshal([]byte(in.Source), &tp); err != nil {
		return &empty.Empty{}, err
	}

	err := s.consumer.StoreOffset(ctx, tp)

	return &empty.Empty{}, err
}

// Quit shuts down the server
func (s *WorkSchedulerServer) Quit(ctx context.Context, in *empty.Empty) (*empty.Empty, error) {
	return &empty.Empty{}, nil
}

// Stop the server from distributing new work
func (s *WorkSchedulerServer) Stop(ctx context.Context, in *empty.Empty) (*empty.Empty, error) {
	// stop the grpcServer
	s.grpcServer.Stop()
	return &empty.Empty{}, nil
}

// Lag reports the lag of the oldest unprocessing message
func (s *WorkSchedulerServer) Lag(ctx context.Context, in *empty.Empty) (*pb.LagResponse, error) {
	return &pb.LagResponse{}, nil
}

// Run the server
func (s *WorkSchedulerServer) Run(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	s.running(ctx)
}

func (s *WorkSchedulerServer) running(ctx context.Context) {
	// send Stop on graceful shutdown
	go func() {
		for {
			select {
			case <-ctx.Done():
				s.Stop(ctx, nil)
			}
		}
	}()

	if err := s.initKafkaConsumer(ctx); err != nil {
		log.Fatalf("failed to initialize kafka consumer: %v", err)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", s.Config.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("starting grpc server on %v", lis.Addr())
	// this blocks until he grpc server exits via a call to s.grpcServer.Stop()
	if err := s.grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	log.Printf("grpc server exited properly")
}

func (s *WorkSchedulerServer) initKafkaConsumer(ctx context.Context) error {
	prefixedTopic := fmt.Sprintf("%s_%s", s.Config.Prefix, s.Config.Topic)
	prefixedGroup := fmt.Sprintf("%s_%s", s.Config.Prefix, s.Config.Group)
	kafkaConsumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": s.Config.Brokers,
		"group": prefixedGroup,
	})
	if err != nil {
		return err
	}

	err = kafkaConsumer.Subscribe(prefixedTopic, nil)
	if err != nil {
		return err
	}

	s.consumer, err = consumer.New(kafkaConsumer)
	return err
}
