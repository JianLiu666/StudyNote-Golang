package cmd

import (
	"fmt"
	"natspractice/config"
	"runtime"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var gitCommitNum string
var buildTime string

var rootCmd = &cobra.Command{
	Use:               "root",
	Short:             "choose which one to run: publisher or subscriber",
	Long:              `This is the test case for NATS streaming.`,
	PersistentPreRunE: PersistentPreRunBeforeCommandStartUp,
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "f", "./conf.d/env.yaml", "config file")
	rootCmd.PersistentFlags().StringVarP(&gitCommitNum, "version", "v", "unknown", "git commit hash")
	rootCmd.PersistentFlags().StringVarP(&buildTime, "buildTime", "b", time.Now().String(), "binary build time")
}

func initConfig() {
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()
	viper.SetEnvPrefix("NATS")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	}

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("ReadInConfig file failed: %v\n", err)
	} else {
		fmt.Printf("Using config file: %v\n", viper.ConfigFileUsed())
	}
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func PersistentPreRunBeforeCommandStartUp(cmd *cobra.Command, args []string) error {
	goVersion := runtime.Version()
	osName := runtime.GOOS
	architecture := runtime.GOARCH
	fmt.Println("======")
	fmt.Printf("Build on %s\n", buildTime)
	fmt.Printf("GoVersion: %s\n", goVersion)
	fmt.Printf("GitCommitNum: %s\n", gitCommitNum)
	fmt.Printf("OS: %s\n", osName)
	fmt.Printf("Architecture: %s\n", architecture)
	fmt.Println("======")

	c, err := config.NewFromViper()
	if err != nil {
		fmt.Printf("Init config failed: %v\n", err)
	} else {
		config.SetConfig(c)
	}

	return nil
}
