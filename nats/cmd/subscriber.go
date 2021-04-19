package cmd

import "github.com/spf13/cobra"

var subscriberCmd = &cobra.Command{
	Use:   "subscriber",
	Short: "Run NATS streaming subscriber test case.",
	Long:  `No more description.`,
	RunE:  RunSubscriberCmd,
}

func init() {
	rootCmd.AddCommand(subscriberCmd)
}

func RunSubscriberCmd(cmd *cobra.Command, args []string) error {

	return nil
}
