package app

import (
	"flag"
	"load-balancer/internal/app/algorithms"
	"load-balancer/internal/app/backend"
	"log"
	"net/url"
	"strconv"
	"strings"
)

type Config struct {
	Port      int
	Backends  []*backend.Backend
	Algorithm algorithms.Algorithm
}

func NewConfig() *Config {
	conf := &Config{
		Backends:  make([]*backend.Backend, 0),
		Algorithm: &algorithms.WeightedRoundRobin{},
	}

	conf.parseFlags()

	return conf
}

func (c *Config) parseFlags() {
	port := flag.Int("port", 8080, "listening port")
	backends := flag.String("backends", "", "comma-separated list of backend addresses")
	backendsWeights := flag.String("backends-weight", "", "comma-separated list of backend's weights")
	flag.Parse()

	if len(*backends) == 0 {
		log.Fatal("at least one backend address is required")
		return
	}

	c.Port = *port

	backendsURLs := strings.Split(*backends, ",")
	backendsWeightsArr := strings.Split(*backendsWeights, ",")
	weights := len(backendsWeightsArr)
	for index, backendURL := range backendsURLs {
		backendURL = strings.TrimSpace(backendURL)
		u, err := url.Parse(backendURL)
		if err != nil {
			log.Fatal(err)
		}

		weight := uint64(1)
		if weights != 0 && index <= weights-1 {
			weight, _ = strconv.ParseUint(backendsWeightsArr[index], 0, 64)
		}

		c.Backends = append(c.Backends, backend.NewBackend(u, uint(weight)))
	}
}
