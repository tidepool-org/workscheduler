// Package main implements a server for Greeter service.
package main

import (
	"context"
	"log"
	"net"

	empty "github.com/golang/protobuf/ptypes/empty"
	pb "github.com/tidepool-org/workscheduler/workscheduler"
	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

// server is used to implement workscheduler
type server struct {
	pb.UnimplementedWorkSchedulerServer
}

func (s *server) Poll(ctx context.Context, in *empty.Empty) (*pb.Work, error) {
	return &pb.Work{}, nil
}

func (s *server) Ping(ctx context.Context, in *empty.Empty) (*empty.Empty, error) {
	return &empty.Empty{}, nil
}

func (s *server) Failed(ctx context.Context, in *pb.WorkSource) (*empty.Empty, error) {
	return &empty.Empty{}, nil
}

func (s *server) Complete(ctx context.Context, in *pb.WorkSource) (*pb.WorkOutput, error) {
	return &pb.WorkOutput{}, nil
}

func (s *server) Quit(ctx context.Context, in *empty.Empty) (*empty.Empty, error) {
	return &empty.Empty{}, nil
}

func (s *server) Stop(ctx context.Context, in *empty.Empty) (*empty.Empty, error) {
	return &empty.Empty{}, nil
}

func (s *server) Stop(ctx context.Context, in *empty.Empty) (*pb.LagResponse, error) {
	return &pb.LagResponse{}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterWorkSchedulerServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
