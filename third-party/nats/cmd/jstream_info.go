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
	Use:   "jsinfo",
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
	for info := range js.StreamsInfo(nats.Context(ctx)) {
		fmt.Println("stream name:", info.Config.Name)
	}

	// Get information about all consumers (with MaxWait JSOpt)
	for info := range js.ConsumersInfo("S", nats.MaxWait(10*time.Second)) {
		fmt.Println("consumer name:", info.Name)
	}

	return nil
}
