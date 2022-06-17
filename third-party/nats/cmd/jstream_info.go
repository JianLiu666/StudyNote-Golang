package cmd

import (
	"context"
	"fmt"
	"jian6/third-party/nats/config"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/spf13/cobra"
)

var jsInfoCmd = &cobra.Command{
	Use:   "js_info",
	Short: "Get NATS JetStream Information.",
	Long:  `No more description.`,
	RunE:  RunJsInfoCmd,
}

func init() {
	rootCmd.AddCommand(jsInfoCmd)
}

func RunJsInfoCmd(cmd *cobra.Command, args []string) error {
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
			fmt.Printf("%s: %s (%s)\n", stream.Config.Name, subject, stream.Config.Retention)
		}
	}

	// Get information about all consumers (with MaxWait JSOpt)
	for info := range js.ConsumersInfo("Collection", nats.MaxWait(10*time.Second)) {
		fmt.Printf("%s: %s\n", info.Stream, info.Name)
	}

	for info := range js.ConsumersInfo("Delivery", nats.MaxWait(10*time.Second)) {
		fmt.Printf("%s: %s\n", info.Stream, info.Name)
	}

	return nil
}
