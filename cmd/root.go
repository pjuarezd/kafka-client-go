package cmd

import (
	"fmt"
	"kafka-client/util"
	"os"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/spf13/cobra"
)

var (
	config  kafka.ConfigMap
	cfgFile string
)

var rootCmd = &cobra.Command{
	Use:   "kafka-client",
	Short: "you know, to talk to kafka",
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(loadConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "client.properties", "broker config")
	rootCmd.AddCommand(producerCdm)
	//rootCmd.AddCommand(consumerCmd)
	rootCmd.AddCommand(topicCmd)
}

func loadConfig() {
	if cfgFile == "" {
		fmt.Println("Unable to open config file, path no provided")
		os.Exit(1)
	}
	if !util.FileExists(cfgFile) {
		fmt.Printf("Unable to open config file, file do not exist: \"%s\"", cfgFile)
		os.Exit(1)
	}
	config = util.ReadConfig(cfgFile)
}