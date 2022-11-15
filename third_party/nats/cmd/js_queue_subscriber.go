package cmd

import (
	"fmt"
	"jian6/third_party/nats/config"
	"os"
	"os/signal"
	"syscall"

	"github.com/nats-io/nats.go"
	"github.com/spf13/cobra"
)

var jsQueueSubscriberCmd = &cobra.Command{
	Use:   "js_qsub",
	Short: "Run NATS JetStream queue group subscriber example.",
	Long:  `No more description.`,
	RunE:  RunJsQueueSubscriberCmd,
}

func init() {
	rootCmd.AddCommand(jsQueueSubscriberCmd)
}

func RunJsQueueSubscriberCmd(cmd *cobra.Command, args []string) error {
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
	_, err = js.QueueSubscribe("Collection.GuChat.Direct", "BLS",
		callback,
		nats.ManualAck(),
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
