package main

import (
	"fmt"
	"math"

	"github.com/sirupsen/logrus"
)

func calculate(clients []*client) {
	var total int64
	var shortest int64 = math.MaxInt64
	var longest int64 = math.MinInt64

	for _, client := range clients {
		for i := 0; i < conf.Simulation.NumMessages; i++ {
			diff := client.times[i]["client_received"] - client.times[i]["client_start"]
			if diff < shortest {
				shortest = diff
			}
			if diff > longest {
				longest = diff
			}
			total += diff
		}
	}

	fmt.Println()
	logrus.Infof("average timestamp: %v ms", total/(int64(conf.Simulation.NumClients*conf.Simulation.NumMessages)))
	logrus.Infof("shortest timestamp: %v ms", shortest)
	logrus.Infof("longest timestamp: %v ms", longest)
}
