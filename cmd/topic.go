package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	client "kafka-client/kafka"
)

var (
	partitions        int
	replicationFactor int
)

func init() {
	topicCmd.Flags().StringVar(&topicName, "topic", "", "name of the topic to create")
	topicCmd.Flags().IntVar(&partitions, "partitions", 2, "number of partitions")
	topicCmd.Flags().IntVar(&replicationFactor, "shards", 1, "replication factor")
	topicCmd.MarkFlagRequired("topic")
}

var topicCmd = &cobra.Command{
	Use:   "topic",
	Short: "Creates a new topic in kafka",
	Run:   runFn,
}

func runFn(cmd *cobra.Command, args []string) {
	fmt.Println("Creating new topic")
	client.CreateTopic(&config, topicName, partitions, replicationFactor)
}
