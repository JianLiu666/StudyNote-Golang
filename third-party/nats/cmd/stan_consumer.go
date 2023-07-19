package cmd

import (
	"fmt"
	"natspractice/config"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/nats-io/stan.go"
	"github.com/spf13/cobra"
)

var stanConsumerCmd = &cobra.Command{
	Use:   "stan_sub",
	Short: "",
	Long:  ``,
	RunE:  RunStanConsumerCmd,
}

func init() {
	rootCmd.AddCommand(stanConsumerCmd)
}

func RunStanConsumerCmd(cmd *cobra.Command, args []string) error {
	sc, err := stan.Connect(
		config.Nats.ClusterId,
		fmt.Sprintf("stna-%v", time.Now().UnixNano()),
		stan.NatsURL(config.Nats.Addr),
	)
	if err != nil {
		return err
	}
	defer sc.Close()

	_, err = sc.Subscribe("Test", func(msg *stan.Msg) {
		eplased, _ := strconv.Atoi(string(msg.Data))
		fmt.Println(time.Now().Sub(time.UnixMilli(int64(eplased))))
	})
	if err != nil {
		return err
	}

	// set graceful shutdown method
	stopSignal := make(chan os.Signal, 1)
	signal.Notify(stopSignal, os.Interrupt, syscall.SIGTERM)
	<-stopSignal

	return nil
}
