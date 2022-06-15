package cmd

import (
	"fmt"
	"jian6/third-party/nats/config"
	"os"
	"os/signal"
	"syscall"

	"github.com/nats-io/nats.go"
	"github.com/spf13/cobra"
)

var jsSubscriberCmd = &cobra.Command{
	Use:   "jssub",
	Short: "Run NATS streaming subscriber test case.",
	Long:  `No more description.`,
	RunE:  RunJsSubscriberCmd,
}

func init() {
	rootCmd.AddCommand(jsSubscriberCmd)
}

func RunJsSubscriberCmd(cmd *cobra.Command, args []string) error {
	nc, err := nats.Connect(config.Nats.Addr, nats.Name("consumer"))
	if err != nil {
		return err
	}
	defer nc.Close()

	js, err := nc.JetStream(nats.PublishAsyncMaxPending(256))
	if err != nil {
		return err
	}

	callback := func(msg *nats.Msg) {
		fmt.Println(string(msg.Data))
		msg.Ack()
	}

	// Create durable consumer monitor
	_, err = js.Subscribe("*", callback,
		nats.ManualAck(),
		nats.Durable("consumer"),
		nats.BindStream("S"))
	if err != nil {
		return err
	}

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM)
	<-stopChan

	return nil
}
