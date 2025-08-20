package main

import (
	"context"
	"errors"
	"fmt"
	"load-balancer/internal/app"
	"load-balancer/internal/app/tools"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	config := app.NewConfig()

	balancer := app.NewLoadBalancer(config.Algorithm, config.Backends)

	setHealthChecker(balancer)

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	server := http.Server{
		Addr:    fmt.Sprintf(":%d", config.Port),
		Handler: balancer,
	}

	go func() {
		log.Printf("Starting balancer on port %d\n", config.Port)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	<-shutdown
	log.Println("Shutting down balancer...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
}

func setHealthChecker(balancer *app.LoadBalancer) {
	healthcheck := &tools.HealthChecker{
		TimeOut:  2 * time.Second,
		Backends: balancer.Backends(),
	}

	healthcheck.Check()

	go healthcheck.CheckWithPeriod(1 * time.Minute)
}
