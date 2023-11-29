package cmd

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var rootCmd = &cobra.Command{
	Use:   "root",
	Short: "",
	Long:  ``,
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVarP(&cfgFile,
		"confg", "f",
		"./config/env.yaml",
		"config file (defaultis ./config/env.yaml)",
	)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logrus.Errorln(err)
		os.Exit(1)
	}
}

func initConfig() {
	viper.SetConfigFile(cfgFile)
	viper.AutomaticEnv()
}
