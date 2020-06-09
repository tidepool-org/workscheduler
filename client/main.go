package main

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	pb "github.com/tidepool-org/workscheduler/workscheduler"
	"google.golang.org/grpc"
	"log"
	"time"
)

const (
	address     = "localhost:5051"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewWorkSchedulerClient(conn)

	work := make([]*pb.Work, 0)
	for i := 0; i < 5; i++ {
		res, err := c.Poll(context.Background(), &empty.Empty{})
		if err != nil {
			log.Printf(err.Error())
			continue
		}
		work = append(work, res)
		log.Printf(string(res.Data))
	}

	log.Printf("%v", work)
	// complete first 2 items in order, next 3 in reverse order
	if _, err := c.Complete(context.Background(), work[0].Source); err != nil {
		log.Printf(err.Error())
	}
	time.Sleep(time.Second)
	if _, err := c.Complete(context.Background(), work[1].Source); err != nil {
		log.Printf(err.Error())
	}
	time.Sleep(time.Second)
	if _, err := c.Complete(context.Background(), work[4].Source); err != nil {
		log.Printf(err.Error())
	}
	time.Sleep(time.Second)
	if _, err := c.Complete(context.Background(), work[3].Source); err != nil {
		log.Printf(err.Error())
	}
	time.Sleep(time.Second)
	if _, err := c.Complete(context.Background(), work[2].Source); err != nil {
		log.Printf(err.Error())
	}
}
