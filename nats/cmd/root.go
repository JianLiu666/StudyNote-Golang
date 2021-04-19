package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

var (
	GitCommitNum string
	BuildTime    string
)

var rootCmd = &cobra.Command{
	Use:               "root",
	Short:             "choose which one to run: publisher or subscriber",
	Long:              `This is the test case for NATS streaming.`,
	PersistentPreRunE: PersistentPreRunBeforeCommandStartUp,
}

func init() {

}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func PersistentPreRunBeforeCommandStartUp(cmd *cobra.Command, args []string) error {
	goVersion := runtime.Version()
	osName := runtime.GOOS
	architecture := runtime.GOARCH
	fmt.Printf("Go Version %s\n", goVersion)
	fmt.Printf("Build on %s from gitCommitNum %s\n", BuildTime, GitCommitNum)
	fmt.Printf("OS: %s, Architecture: %s\n", osName, architecture)

	return nil
}
