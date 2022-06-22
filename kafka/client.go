package client

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func CreateTopic(config *kafka.ConfigMap, topicName string, part int, rf int) {
	adm, err := kafka.NewAdminClient(config)
	if err != nil {
		fmt.Printf("Failed to conect to broker: %s/n", err)
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	maxTime, errt := time.ParseDuration("60s")
	if err != nil {
		panic("ParseDuration(60s)")

	}
	specification := []kafka.TopicSpecification{
		{
			Topic:             topicName,
			NumPartitions:     part,
			ReplicationFactor: rf}}

	topicResult, errdt := adm.CreateTopics(ctx, specification, kafka.SetAdminOperationTimeout(maxTime))

	if errdt != nil {
		fmt.Printf("Failed to create topic")
	}

	if errt != nil {
		fmt.Printf("Failed to Create topic \"%s\": %s", topicName, errdt)
		os.Exit(1)
	}

	for _, result := range topicResult {
		fmt.Printf("%s\n", result)
	}

	adm.Close()
}

func CreateConsumerAndSuscribe(config *kafka.ConfigMap, topic string) *kafka.Consumer {

	cns, err := kafka.NewConsumer(config)

	if err != nil {
		fmt.Printf("Error creating consumer: %s", err)
	}

	cns.SubscribeTopics([]string{topic}, nil)

	return cns
}

func CreateProducer(config *kafka.ConfigMap) *kafka.Producer {
	prd, err := kafka.NewProducer(config)

	if err != nil {
		fmt.Printf("Error creating Producer: %s", err)
		os.Exit(1)
	}

	return prd
}
