package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/spf13/cobra"
	client "kafka-client/kafka"
)

func init() {
	consumerCmd.Flags().StringVar(&topicName, "topic", "", "name of the topic to listen")
	consumerCmd.MarkFlagRequired("topic")
}

var consumerCmd = &cobra.Command{
	Use:   "consumer",
	Short: "Reads from a kafka topic and prints on screen until finishes",
	Run:   consumerCmdFn,
}

func consumerCmdFn(cmd *cobra.Command, args []string) {
	fmt.Println("Creating Producer")
	kfkConsumer := client.CreateConsumerAndSuscribe(&config, topicName)
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	run := true
	for run {
		select {
		case sig := <-sigchan:
			fmt.Printf("Caught signal %v: terminating\n", sig)
			run = false
		default:
			ev := kfkConsumer.Poll(100)
			if ev == nil {
				continue
			}

			switch e := ev.(type) {
			case *kafka.Message:
				fmt.Printf("%% Message on %s:\n%s\n",
					e.TopicPartition, string(e.Value))
				if e.Headers != nil {
					fmt.Printf("%% Headers: %v\n", e.Headers)
				}
				_, err := kfkConsumer.StoreMessage(e)
				if err != nil {
					fmt.Fprintf(os.Stderr, "%% Error storing offset after message %s:\n",
						e.TopicPartition)
				}
			case kafka.Error:
				fmt.Fprintf(os.Stderr, "%% Error: %v: %v\n", e.Code(), e)
				if e.Code() == kafka.ErrAllBrokersDown {
					run = false
				}
			default:
				fmt.Printf("Ignored %v\n", e)
			}
		}
	}

	fmt.Printf("Closing consumer\n")
	kfkConsumer.Close()
}
