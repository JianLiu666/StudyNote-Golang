package cmd

import (
	"fmt"
	"jian6/third-party/nats/config"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/spf13/cobra"
)

var jsPublisherCmd = &cobra.Command{
	Use:   "jspub",
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
		return err
	}
	defer nc.Close()

	// Create JetStream Context
	js, err := nc.JetStream(nats.PublishAsyncMaxPending(256))
	if err != nil {
		return err
	}

	js.Publish("SS.subj1", []byte(fmt.Sprintf("%d", time.Now().UnixMilli())))

	select {
	case <-js.PublishAsyncComplete():
	}

	fmt.Println("Published.")
	return nil
}
