package cmd

import (
	"jian6/third-party/nats/config"
	"log"

	"github.com/nats-io/nats.go"
	"github.com/spf13/cobra"
)

var jsAddStreamCmd = &cobra.Command{
	Use:   "js_add_stream",
	Short: "Add new stream to NATS JetStream",
	Long:  `No more description.`,
	RunE:  RunJsAddStreamCmd,
}

func init() {
	rootCmd.AddCommand(jsAddStreamCmd)
}

func RunJsAddStreamCmd(cmd *cobra.Command, args []string) error {
	nc, err := nats.Connect(config.Nats.Addr)
	if err != nil {
		log.Println(err)
		return nil
	}
	defer nc.Close()

	js, err := nc.JetStream()
	if err != nil {
		log.Println(err)
		return nil
	}

	// Create a stream
	_, err = js.AddStream(&nats.StreamConfig{
		Name:      "GuChat_Collection",
		Subjects:  []string{"Collection.GuChat.Direct", "Collection.GuChat.Group"},
		Replicas:  3,
		Retention: nats.WorkQueuePolicy,
	})
	if err != nil {
		log.Println(err)
		return nil
	}

	_, err = js.AddStream(&nats.StreamConfig{
		Name:      "GuChat_Delivery",
		Subjects:  []string{"Delivery.GuChat.Direct", "Delivery.GuChat.Group"},
		Replicas:  3,
		Retention: nats.InterestPolicy,
	})
	if err != nil {
		log.Println(err)
		return nil
	}

	_, err = js.AddStream(&nats.StreamConfig{
		Name:      "KKGame_Collection",
		Subjects:  []string{"Collection.KKGame.Group"},
		Replicas:  3,
		Retention: nats.WorkQueuePolicy,
	})
	if err != nil {
		log.Println(err)
		return nil
	}

	_, err = js.AddStream(&nats.StreamConfig{
		Name:      "KKGame_Delivery",
		Subjects:  []string{"Delivery.KKGame.Group"},
		Replicas:  3,
		Retention: nats.InterestPolicy,
	})
	if err != nil {
		log.Println(err)
		return nil
	}

	return nil
}
