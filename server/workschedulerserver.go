package main

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes/empty"
	pb "github.com/tidepool-org/workscheduler/workscheduler"
	"google.golang.org/grpc"
	"log"
	"net"
	"sync"
)

// Config is the configuration of the work scheduler
type Config struct {
	brokers string
		prefix  string
	port    string
}

// WorkSchedulerServer is used to implement workscheduler
type WorkSchedulerServer struct {
	pb.UnimplementedWorkSchedulerServer
	Config     Config
	grpcServer *grpc.Server
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
	return &pb.Work{}, nil
}

// Ping to check for health of work scheduler server
func (s *WorkSchedulerServer) Ping(ctx context.Context, in *empty.Empty) (*empty.Empty, error) {
	return &empty.Empty{}, nil
}

// Failed to report failed job
func (s *WorkSchedulerServer) Failed(ctx context.Context, in *pb.WorkSource) (*empty.Empty, error) {
	return &empty.Empty{}, nil
}

// Complete to report completed job
func (s *WorkSchedulerServer) Complete(ctx context.Context, in *pb.WorkSource) (*pb.WorkOutput, error) {
	return &pb.WorkOutput{}, nil
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

	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", s.Config.port))
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
