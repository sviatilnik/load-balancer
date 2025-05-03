package main

import (
	"fmt"
	"load-balancer/internal/app"
	"load-balancer/internal/app/tools"
	"log"
	"net/http"
	"time"
)

func main() {
	config := app.NewConfig()

	balancer := app.NewLoadBalancer(config.Algorithm, config.Backends)

	setHealthChecker(balancer)

	log.Printf("Load Balancer starting at: %d\n", config.Port)

	if err := http.ListenAndServe(fmt.Sprintf(":%d", config.Port), balancer); err != nil {
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
