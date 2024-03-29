package cmd

import (
	"fmt"
	"natspractice/config"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/spf13/cobra"
)

var jsPublisherCmd = &cobra.Command{
	Use:   "js_pub",
	Short: "A nats jetstream example: publish message",
	Long:  ``,
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
	subjNames := []string{
		"Collection.GuChat.Direct",
		// "Collection.GuChat.Group",
		// "Collection.KKGame.Group",
		// "Collection.KKGame.Group",
		// "Delivery.GuChat.Direct",
		// "Delivery.KKGame.Group",
	}
	for _, subjName := range subjNames {
		for i := 0; i < 1000; i++ {
			msg := fmt.Sprintf("%d:%04d", time.Now().UnixMilli(), i)
			_, err := js.Publish(subjName, []byte(msg))
			if err != nil {
				fmt.Printf("%v: %v", subjName, err)
				return nil
			}
		}
	}

	select {
	case <-js.PublishAsyncComplete():
	}

	fmt.Println("Published.")
	return nil
}
