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

var subscriberCmd = &cobra.Command{
	Use:   "sub",
	Short: "Run NATS streaming subscriber test case.",
	Long:  `No more description.`,
	RunE:  RunSubscriberCmd,
}

func init() {
	rootCmd.AddCommand(subscriberCmd)
}

func RunSubscriberCmd(cmd *cobra.Command, args []string) error {
	opts := []nats.Option{nats.Name("Client-Consumer")}

	nc, err := nats.Connect(config.Nats.Addr, opts...)
	if err != nil {
		return err
	}
	defer nc.Close()

	count := 0
	_, err = nc.Subscribe("learning_note", func(msg *nats.Msg) {
		count++
		fmt.Printf("[#%d] Received on [learning_note]: '%s'\n", count, string(msg.Data))
	})
	if err != nil {
		return err
	}
	nc.Flush()

	if err := nc.LastError(); err != nil {
		return err
	}

	fmt.Println("Listening on [learning_note]")

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM)
	<-stopChan

	return nil
}
