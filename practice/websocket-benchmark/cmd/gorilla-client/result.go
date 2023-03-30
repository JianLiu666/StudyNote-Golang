package main

import (
	"fmt"
	"math"

	"github.com/sirupsen/logrus"
)

func calculate(clients [numClients]*client) {
	var total int64
	var shortest int64 = math.MaxInt64
	var longest int64 = math.MinInt64

	for i := 0; i < numClients; i++ {
		for j := 1; j <= numMessages; j++ {
			diff := clients[i].times[int32(j)]["client_received"] - clients[i].times[int32(j)]["client_start"]
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
	logrus.Infof("average timestamp: %v ms", total/(int64(numClients)*int64(numMessages)))
	logrus.Infof("shortest timestamp: %v ms", shortest)
	logrus.Infof("longest timestamp: %v ms", longest)
}
