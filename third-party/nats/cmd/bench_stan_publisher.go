package cmd

import (
	"fmt"
	"natspractice/config"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nats-io/stan.go"
	"github.com/spf13/cobra"
)

var benchStanPublisherCmd = &cobra.Command{
	Use:   "bench_stan_pub",
	Short: "",
	Long:  ``,
	RunE:  RunBenchStanPublisherCmd,
}

func init() {
	rootCmd.AddCommand(benchStanPublisherCmd)
}

func RunBenchStanPublisherCmd(cmd *cobra.Command, args []string) error {
	sc, err := stan.Connect(
		config.Nats.StanClusterId,
		fmt.Sprintf("stan-%v", time.Now().UnixNano()),
		stan.NatsURL(config.Nats.Addr),
	)
	if err != nil {
		return err
	}
	defer sc.Close()

	// 併發測試 NATS streaming publish 效能
	for i := 0; i < config.Nats.BenchNumProducers; i++ {
		go func() {
			for i := 0; i < config.Nats.BenchProducerEachTimes; i++ {
				eplased := time.Now()
				err = sc.Publish("Test", []byte(fmt.Sprintf("%v", time.Now().UnixMilli())))
				fmt.Println(i, time.Now().Sub(eplased))

				if err != nil {
					fmt.Println(err)
				}

				time.Sleep(100 * time.Millisecond)
			}
		}()
	}

	// set graceful shutdown method
	stopSignal := make(chan os.Signal, 1)
	signal.Notify(stopSignal, os.Interrupt, syscall.SIGTERM)
	<-stopSignal

	return nil
}
