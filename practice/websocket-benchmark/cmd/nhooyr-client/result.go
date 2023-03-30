package main

import (
	"fmt"
	"sort"

	"github.com/sirupsen/logrus"
)

func calculate(clients []*client) {
	var total int64

	data := []int64{}

	for _, client := range clients {
		for i := 0; i < conf.Simulation.NumMessages; i++ {
			diff := client.times[i]["client_received"] - client.times[i]["client_start"]
			total += diff
			data = append(data, diff)
		}
	}

	sort.Slice(data, func(i, j int) bool {
		return data[i] < data[j]
	})

	fmt.Println()
	logrus.Info("timestmap")
	logrus.Infof("average:  %v ms", total/(int64(conf.Simulation.NumClients*conf.Simulation.NumMessages)))
	logrus.Infof("shortest: %v ms", data[0])
	logrus.Infof("longest:  %v ms", data[len(data)-1])
	logrus.Infof("p50:      %v ms", data[int(float64(len(data))*0.5)])
	logrus.Infof("p90:      %v ms", data[int(float64(len(data))*0.9)])
	logrus.Infof("p95:      %v ms", data[int(float64(len(data))*0.95)])
}
