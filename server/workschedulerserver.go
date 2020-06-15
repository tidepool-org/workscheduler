package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/tidepool-org/workscheduler/orchestrator"
	pb "github.com/tidepool-org/workscheduler/workscheduler"
	"google.golang.org/grpc"
	"log"
	"net"
	"sync"
	"time"
)

// Config is the configuration of the work scheduler
type Config struct {
	Brokers     string
	Group       string
	Prefix      string
	Topic       string
	Port        string
	WorkTimeout time.Duration
}

// WorkSchedulerServer is used to implement workscheduler
type WorkSchedulerServer struct {
	pb.UnimplementedWorkSchedulerServer
	Config           Config
	grpcServer       *grpc.Server
	workOrchestrator orchestrator.Orchestrator
}

// NewWorkSchedulerServer create a new WorkSchedulerServer
func NewWorkSchedulerServer(c Config) *WorkSchedulerServer {
	grpcServer := grpc.NewServer()
	wsServer := &WorkSchedulerServer{
		Config:     c,
		grpcServer: grpcServer,
	}

	pb.RegisterWorkSchedulerServer(grpcServer, wsServer)
	return wsServer
}

// Poll for work to process
func (s *WorkSchedulerServer) Poll(ctx context.Context, in *empty.Empty) (*pb.Work, error) {
	select {
	case <-ctx.Done():
		return nil, errors.New("timed out while polling for work")
	case work := <-s.workOrchestrator.WorkChannel():
		return &work, nil
	}
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
	s.workOrchestrator.ResponseChannel() <- *in
	return &empty.Empty{}, nil
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
	s.running(ctx, wg)
}

func (s *WorkSchedulerServer) running(ctx context.Context, wg *sync.WaitGroup) {
	// send Stop on graceful shutdown
	go func() {
		for {
			select {
			case <-ctx.Done():
				s.Stop(ctx, nil)
			}
		}
	}()

	if err := s.startOrchestrator(ctx, wg); err != nil {
		log.Fatalf("Failed to start kafka orchestrator: %v", err)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%v", s.Config.Port))
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

func (s *WorkSchedulerServer) startOrchestrator(ctx context.Context, wg *sync.WaitGroup) error {
	params := orchestrator.KafkaWorkOrchestratorParams{
		Brokers:     s.Config.Brokers,
		Prefix:      s.Config.Prefix,
		Topic:       s.Config.Topic,
		WorkTimeout: s.Config.WorkTimeout,
	}
	wo, err := orchestrator.NewKafkaWorkOrchestrator(params)
	if err != nil {
		return err
	}

	s.workOrchestrator = wo
	wg.Add(1)
	go s.workOrchestrator.Run(ctx, wg)
	return nil
}
