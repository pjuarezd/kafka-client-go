package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	client "kafka-client/kafka"
	util "kafka-client/util"
)

func init() {
	rootCmd.AddCommand(topicCmd)
}

var topicCmd = &cobra.Command{
	Use: "topic",
	Short: "Creates a new topic in kafka",
	Run: runFn}

func runFn(cmd *cobra.Command, args []string) {
	conf := util.ReadConfig("client.properties")
	fmt.Println("Creating new topic")
	client.CreateTopic(&conf, "", 2, 1)
}