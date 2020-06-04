package main

import (
	"context"
	"log"
	"os"
	"time"

	"google.golang.org/grpc"
	empty "github.com/golang/protobuf/ptypes/empty"
	pb "github.com/tidepool-org/workscheduler/workscheduler"
)

const (
	address     = "localhost:50051"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewWorkSchedulerClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.Ping(ctx, &empty.Emptuy{})
	if err != nil {
		log.Fatalf("could not ping: %v", err)
	}
}
