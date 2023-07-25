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

var benchJetStreamPublisherCmd = &cobra.Command{
	Use:   "bench_js_pub",
	Short: "A simple nats jetstream benchmark example, cooperate with 'bench_js_sub' command.",
	Long:  ``,
	RunE:  RunBenchJetStreamPublisherCmd,
}

func init() {
	rootCmd.AddCommand(benchJetStreamPublisherCmd)
}

func RunBenchJetStreamPublisherCmd(cmd *cobra.Command, args []string) error {
	nc, err := nats.Connect(config.Nats.Addr)
	if err != nil {
		return err
	}
	defer nc.Close()

	js, err := nc.JetStream(nats.PublishAsyncMaxPending(256))
	if err != nil {
		return err
	}

	// Check if the stream already exists
	_, err = js.StreamInfo("Test")
	if err != nil {
		if err == nats.ErrStreamNotFound {
			// If the stream does not exist, create it
			_, err = js.AddStream(&nats.StreamConfig{
				Name:     "Test",
				Subjects: []string{"Test"},
				Replicas: 1,
				Storage:  nats.FileStorage,
			})

			if err != nil {
				return err
			}
		} else {
			// If there was another error, return it
			return err
		}
	}

	// 併發測試 NATS JetStream publish 效能
	for i := 0; i < config.Nats.BenchNumProducers; i++ {
		go func() {
			for i := 0; i < config.Nats.BenchProducerEachTimes; i++ {
				eplased := time.Now()
				_, err = js.Publish("Test", []byte(fmt.Sprintf("%v", time.Now().UnixMilli())))
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
