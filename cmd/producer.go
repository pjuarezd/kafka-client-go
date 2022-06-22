package cmd

import (
	"fmt"
	"kafka-client/util"
	"os"

	client "kafka-client/kafka"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/spf13/cobra"
)

var (
	topicName   string
	message     string
	headers     string
	delay       int
	csvFilePath string
)

type KafkaMessage struct {
	Message []byte
	Headers []kafka.Header
}

func init() {
	producerCdm.Flags().StringVarP(&topicName, "topic", "r", "", "name of the topic to listen")
	producerCdm.Flags().StringVar(&message, "message", "", "message to send to kafka")
	producerCdm.Flags().StringVar(&headers, "headers", "", "Key value, comma separated to add as headers of the message")
	producerCdm.Flags().IntVar(&delay, "delay", 1, "Delay in seconds between messages")
	producerCdm.Flags().StringVar(&csvFilePath, "file", "", "name of the topic to listen")
	producerCdm.MarkFlagRequired("topic")
	producerCdm.MarkFlagsRequiredTogether("message", "headers")
	producerCdm.MarkFlagsMutuallyExclusive("message", "file")
	producerCdm.MarkFlagsMutuallyExclusive("headers", "file")
}

var producerCdm = &cobra.Command{
	Use:   "producer",
	Short: "Produces a series of messages to the Kafka broker",
	Long: `

Produces a messages to the Kafka broker, or a series of messages from a CSV formatted file.
Either use the "message" and "headers" parameters, or the "file" parameter. 
If you provider the "file" parameter the "message" and "headers" parameters are ignored.

CSV File
=========
When a CSV file is provided, it expects a line with the field names, ie:
-----------------
vegetables,fruits
carrot,banana
potato,strawberry
-----------------
In the example, the first line is the Headers.

The first column will be sent as the message, as of the rest of the lines will be sent to kafka as "headers"
`,
	Run: producerCmdFn,
}

func producerCmdFn(cmd *cobra.Command, args []string) {
	fmt.Println("Creating Producer")
	kfkProducer := client.CreateProducer(&config)

	go func() {
		for e := range kfkProducer.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				m := ev
				if m.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", m.TopicPartition.Error)
				} else {
					fmt.Printf("Delivered message to topic %s [%d] at offset %v\n",
						*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
				}
			case kafka.Error:
				fmt.Printf("Error: %v\n", ev)
			default:
				fmt.Printf("Ignored event: %s\n", ev)
			}
		}
	}()

	if !util.FileExists(csvFilePath) {
		fmt.Printf("Unable to open config file, file do not exist: %s", cfgFile)
		os.Exit(1)
	}

	messages := ParseKafmaMessagesFromCsv(csvFilePath)

	for i, msg := range messages {
		err := SendMessage(kfkProducer, msg)
		if err != nil {
			fmt.Printf("Problem sending message %d: %s", i, err)
			continue
		}
	}
	for kfkProducer.Flush(10000) > 0{
		fmt.Println("Still waiting to flush outstanding messages")
	}
	defer kfkProducer.Close()
}

func ParseKafmaMessagesFromCsv(cvsFilePath string) []KafkaMessage {
	var messages []KafkaMessage
	var headers []string
	for index, line := range util.ReadCSV(cvsFilePath) {
		if index == 0 {
			headers = line
			continue
		}
		var kHeaders []kafka.Header

		for fi, field := range line {
			header := kafka.Header{Key: headers[fi], Value: []byte(field)}
			kHeaders = append(kHeaders, header)
		}
		messages = append(messages, KafkaMessage{Message: []byte(line[0]), Headers: []kafka.Header{}})
	}
	return messages
}

func SendMessage(producer *kafka.Producer, message KafkaMessage) error {
	err := producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topicName, Partition: kafka.PartitionAny},
		Value:          []byte(message.Message),
		Headers:        message.Headers,
	}, nil)
	return err
}
