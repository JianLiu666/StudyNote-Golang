package cmd

import (
	"fmt"
	"natspractice/config"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/spf13/cobra"
)

var benchNatsPublisherCmd = &cobra.Command{
	Use:   "bench_nats_pub",
	Short: "A simple nats benchmark example, cooperate with 'bench_nats_sub' command.",
	Long:  ``,
	RunE:  RunBenchNatsPublisherCmd,
}

func init() {
	rootCmd.AddCommand(benchNatsPublisherCmd)
}

func RunBenchNatsPublisherCmd(cmd *cobra.Command, args []string) error {
	nc, err := nats.Connect(
		config.Nats.Addr,
	)
	if err != nil {
		return err
	}
	defer nc.Close()

	// 併發測試 NATS streaming publish 效能
	for i := 0; i < config.Nats.BenchNumProducers; i++ {
		go func() {
			for i := 0; i < config.Nats.BenchProducerEachTimes; i++ {
				eplased := time.Now()
				err = nc.Publish("Test", []byte(fmt.Sprintf("%v", time.Now().UnixMilli())))
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
