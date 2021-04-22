package cmd

import (
	"fmt"
	"jian6/nats/config"
	"sync"

	"github.com/nats-io/nats.go"
	"github.com/spf13/cobra"
)

var requestReplyCmd = &cobra.Command{
	Use:   "req",
	Short: "Run NATS client for request-reply test case.",
	Long:  `No more description.`,
	RunE:  RunRequestReplyCmd,
}

func init() {
	rootCmd.AddCommand(requestReplyCmd)
}

func RunRequestReplyCmd(cmd *cobra.Command, args []string) error {
	opts := []nats.Option{nats.Name("Client-For-Request")}
	nc, err := nats.Connect(config.Nats.Addr, opts...)
	if err != nil {
		return err
	}
	defer nc.Close()

	wg := sync.WaitGroup{}
	wg.Add(1)

	// Step.1 Prepare a reply-subject and subscribe it
	repSubject := nats.NewInbox()
	_, err = nc.Subscribe(repSubject, func(msg *nats.Msg) {
		fmt.Println(string(msg.Data))
		wg.Done()
	})
	if err != nil {
		return err
	}

	// Step.2 Subscribe target subject
	_, err = nc.Subscribe("case_request_reply", func(msg *nats.Msg) {
		fmt.Println(string(msg.Data))
		if err := msg.Respond([]byte("pong")); err != nil {
			fmt.Println(err)
		}
	})
	if err != nil {
		return err
	}

	// Step.3 Publish a request to target subject
	if err := nc.PublishRequest("case_request_reply", repSubject, []byte("ping")); err != nil {
		return err
	}
	nc.Flush()

	wg.Wait()
	return nil
}
