package cmd

import (
	"fmt"
	"jian6/nats/config"

	"github.com/nats-io/nats.go"
	"github.com/spf13/cobra"
)

var publisherCmd = &cobra.Command{
	Use:   "pub",
	Short: "Run NATS streaming publisher test case.",
	Long:  `No more description.`,
	RunE:  RunPublisherCmd,
}

func init() {
	rootCmd.AddCommand(publisherCmd)
}

func RunPublisherCmd(cmd *cobra.Command, args []string) error {
	opts := []nats.Option{nats.Name("Client-Producer")}

	nc, err := nats.Connect(config.Nats.Addr, opts...)
	if err != nil {
		return err
	}
	defer nc.Close()

	if err := nc.Publish("learning_note", []byte("something")); err != nil {
		return err
	}
	nc.Flush()

	if err := nc.LastError(); err != nil {
		return err
	}
	fmt.Println("Published.")

	return nil
}
