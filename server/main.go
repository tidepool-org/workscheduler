// Package main implements a server for Greeter service.
package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	empty "github.com/golang/protobuf/ptypes/empty"
	pb "github.com/tidepool-org/workscheduler/workscheduler"
	"google.golang.org/grpc"
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
	Config Config
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
	return &empty.Empty{}, nil
}

// Lag reports the lag of the oldest unprocessing message
func (s *WorkSchedulerServer) Lag(ctx context.Context, in *empty.Empty) (*pb.LagResponse, error) {
	return &pb.LagResponse{}, nil
}

// Run the server
func (s *WorkSchedulerServer) Run(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	lis, err := net.Listen("tcp", s.Config.port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterWorkSchedulerServer(grpcServer, &WorkSchedulerServer{})

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func main() {
	var config Config
	var found bool

	config.brokers, found = os.LookupEnv("KAFKA_BROKERS")
	if !found {
		panic("kafka brokers not provided")
	}

	config.prefix, found = os.LookupEnv("KAFKA_TOPIC_PREFIX")
	if !found {
		panic("kafka prefix not provided")
	}

	config.port, found = os.LookupEnv("KAFKA_PORT")
	if !found {
		panic("kafka port not provided")
	}

	server := WorkSchedulerServer{Config: config}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancelFunc := context.WithCancel(context.Background())

	var wg sync.WaitGroup
	wg.Add(1)
	server.Run(ctx, &wg)

	go func(stop chan os.Signal, cancelFunc context.CancelFunc) {
		<-stop
		log.Print("sigint or sigterm received!!!")
		cancelFunc()
	}(stop, cancelFunc)

	wg.Wait()
}
