package cmd

import (
	"context"
	"fmt"
	"jian6/third_party/nats/config"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/spf13/cobra"
)

var jsAddConsumerCmd = &cobra.Command{
	Use:   "js_add_consumer",
	Short: "Add new consumer to NATS JetStream",
	Long:  `No more description.`,
	RunE:  RunJsAddConsumerCmd,
}

func init() {
	rootCmd.AddCommand(jsAddConsumerCmd)
}

func RunJsAddConsumerCmd(cmd *cobra.Command, args []string) error {
	nc, err := nats.Connect(config.Nats.Addr)
	if err != nil {
		return err
	}
	defer nc.Close()

	js, err := nc.JetStream(nats.PublishAsyncMaxPending(256))
	if err != nil {
		return err
	}

	// Get information about all streams (with Context JSOpt)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	for stream := range js.StreamsInfo(nats.Context(ctx)) {
		for _, subject := range stream.Config.Subjects {
			fmt.Printf("%s:%s (%s)\n", stream.Config.Name, subject, stream.Config.Retention)
		}
	}

	// Get information about all consumers (with MaxWait JSOpt)
	for info := range js.ConsumersInfo("S", nats.MaxWait(10*time.Second)) {
		fmt.Printf("%s:%s\n", info.Stream, info.Name)
	}

	return nil
}
