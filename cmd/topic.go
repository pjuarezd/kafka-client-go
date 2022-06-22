package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	client "kafka-client/kafka"
)

func init() {
	topicCmd.PersistentFlags().StringVar(&topicName, "topic", "", "name of the topic to listen")
}

var topicCmd = &cobra.Command{
	Use: "topic",
	Short: "Creates a new topic in kafka",
	Run: runFn,
}

func runFn(cmd *cobra.Command, args []string) {
	fmt.Println("Creating new topic")
	client.CreateTopic(&config, topicName, 2, 1)
}
