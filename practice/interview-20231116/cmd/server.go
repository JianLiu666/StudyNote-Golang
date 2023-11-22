package cmd

import (
	"context"
	"interview20231116/api"
	"interview20231116/pkg/accessor"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "",
	Long:  ``,
	RunE:  RunServerCmd,
}

func init() {
	rootCmd.AddCommand(serverCmd)
}

func RunServerCmd(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	infra := accessor.Build()
	defer infra.Close(ctx)

	infra.InitKvStore(ctx)

	app := api.Init(infra.Config.Server)
	defer app.Shutdown()

	app.Run()

	stopSignal := make(chan os.Signal, 1)
	signal.Notify(stopSignal, os.Interrupt, syscall.SIGTERM)
	<-stopSignal

	logrus.Infof("main: %s closed.\n", cmd.Name())

	return nil
}
