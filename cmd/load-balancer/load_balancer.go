package main

import (
	"flag"
	"fmt"
	"load-balancer/internal/app"
	"load-balancer/internal/app/algorithms"
	"load-balancer/internal/app/backend"
	"load-balancer/internal/app/tools"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
)

func main() {
	port := flag.Int("port", 8080, "listening port")
	backends := flag.String("backends", "", "comma-separated list of backend addresses")
	flag.Parse()

	if len(*backends) == 0 {
		log.Println("at least one backend address is required")
		return
	}

	backendsList := strings.Split(*backends, ",")

	balancer := app.NewLoadBalancer(&algorithms.RoundRobin{})

	for _, back := range backendsList {
		u, err := url.Parse(back)
		if err != nil {
			log.Fatal(err)
		}

		proxy := httputil.NewSingleHostReverseProxy(u)
		proxy.ErrorHandler = func(writer http.ResponseWriter, request *http.Request, e error) {
			log.Println(e.Error())
			http.Error(writer, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
		}

		balancer.AddBackend(backend.NewBackend(u, proxy))
	}

	healthcheck := &tools.HealthChecker{
		TimeOut:  2 * time.Second,
		Backends: balancer.Backends(),
	}

	healthcheck.Check()

	go healthcheck.CheckWithPeriod(1 * time.Minute)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", *port),
		Handler: balancer,
	}

	log.Printf("Load Balancer started at: %d\n", *port)

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
