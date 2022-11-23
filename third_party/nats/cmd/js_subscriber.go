package cmd

import (
	"fmt"
	"natspractice/config"
	"os"
	"os/signal"
	"syscall"

	"github.com/nats-io/nats.go"
	"github.com/spf13/cobra"
)

var jsSubscriberCmd = &cobra.Command{
	Use:   "js_sub",
	Short: "Run NATS JetStream subscriber example.",
	Long:  `No more description.`,
	RunE:  RunJsSubscriberCmd,
}

func init() {
	rootCmd.AddCommand(jsSubscriberCmd)
}

func RunJsSubscriberCmd(cmd *cobra.Command, args []string) error {
	nc, err := nats.Connect(config.Nats.Addr)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer nc.Close()

	js, err := nc.JetStream(nats.PublishAsyncMaxPending(256))
	if err != nil {
		fmt.Println(err)
		return nil
	}

	callback := func(msg *nats.Msg) {
		fmt.Println(string(msg.Data))
		msg.Ack()
	}

	// Create durable consumer monitor
	_, err = js.Subscribe("Collection.GuChat.Direct", callback,
		nats.ManualAck(),
		nats.Durable("consumer1"),
		nats.BindStream("GuChat_Collection"),
	)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM)
	<-stopChan

	return nil
}
