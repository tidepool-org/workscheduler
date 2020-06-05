// Package main implements a server for Greeter service.
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	)


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

	config.port, found = os.LookupEnv("WORK_SCHEDULER_PORT")
	if !found {
		panic("kafka port not provided")
	}

	workSchedulerServer := NewWorkSchedulerServer(config)

	prometheusPort, found := os.LookupEnv("PROMETHEUS_SERVER_PORT")
	if !found {
		panic("prometheus server port not provided")
	}

	promServer := NewPromServer(prometheusPort)

	// listen to signals to stop server
	// convert to cancel on context that server listens to
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	ctx, cancelFunc := context.WithCancel(context.Background())
	go func(stop chan os.Signal, cancelFunc context.CancelFunc) {
		<-stop
		log.Print("sigint or sigterm received!!!")
		cancelFunc()
	}(stop, cancelFunc)

	var wg sync.WaitGroup
	wg.Add(2)
	go workSchedulerServer.Run(ctx, &wg)
	go promServer.Run(ctx, &wg)

	wg.Wait()
}
