package cmd

import "github.com/spf13/cobra"

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

	return nil
}
