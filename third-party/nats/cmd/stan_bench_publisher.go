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

var stanPublisherCmd = &cobra.Command{
	Use:   "stan_pub",
	Short: "",
	Long:  ``,
	RunE:  RunStanPublisherCmd,
}

func init() {
	rootCmd.AddCommand(stanPublisherCmd)
}

func RunStanPublisherCmd(cmd *cobra.Command, args []string) error {
	sc, err := stan.Connect(
		config.Nats.ClusterId,
		fmt.Sprintf("stan-%v", time.Now().UnixNano()),
		stan.NatsURL(config.Nats.Addr),
	)
	if err != nil {
		return err
	}
	defer sc.Close()

	// 併發測試 NATS streaming publish 效能
	for i := 0; i < 10000; i++ {
		go func() {
			for {
				eplased := time.Now()
				err = sc.Publish("Test", []byte(fmt.Sprintf("%v", time.Now().UnixMilli())))
				fmt.Println(time.Now().Sub(eplased))

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
