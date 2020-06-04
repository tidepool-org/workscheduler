package workscheduler

import (
	"context"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type PromServer struct {
	srv *http.Server
}

func NewPromServer(port string) *PromServer {
	http.Handle("/metrics", promhttp.Handler())
	return &PromServer{
		srv: &http.Server{
			Addr: port,
		},
	}
}

func (p *PromServer) Run(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	p.running(ctx)
}

func (p *PromServer) running(ctx context.Context) {
	go func() {
		if err := p.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen:%+s\n", err)
		}
	}()

	<-ctx.Done()

	ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	var err error
	if err = p.srv.Shutdown(ctxShutDown); err != nil {
		log.Fatalf("server Shutdown Failed:%+s", err)
	}

	log.Printf("server exited properly")
}
