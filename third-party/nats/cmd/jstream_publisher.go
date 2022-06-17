package cmd

import (
	"fmt"
	"jian6/third-party/nats/config"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/spf13/cobra"
)

var jsPublisherCmd = &cobra.Command{
	Use:   "js_pub",
	Short: "Run NATS JetStream publisher test case.",
	Long:  `No more description.`,
	RunE:  RunJsPublisherCmd,
}

func init() {
	rootCmd.AddCommand(jsPublisherCmd)
}

func RunJsPublisherCmd(cmd *cobra.Command, args []string) error {
	// Connect to NATS
	nc, err := nats.Connect(config.Nats.Addr)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer nc.Close()

	// Create JetStream Context
	js, err := nc.JetStream(nats.PublishAsyncMaxPending(256))
	if err != nil {
		fmt.Println(err)
		return nil
	}

	// Simple Async Stream Publisher
	subjName := "Collection.GuChat.Direct"
	for i := 0; i < 500; i++ {
		msg := fmt.Sprintf("%d", time.Now().UnixMilli())
		_, err := js.Publish(subjName, []byte(msg))
		if err != nil {
			fmt.Println(err)
			return nil
		}
	}

	select {
	case <-js.PublishAsyncComplete():
	}

	fmt.Println("Published.")
	return nil
}
