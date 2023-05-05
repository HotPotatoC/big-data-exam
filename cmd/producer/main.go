package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/HotPotatoC/bigdata-exam/model"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

func readJSONFile(filename string) ([]byte, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func produceMessage(p *kafka.Producer, topic string, message any) error {

	messageBytes, err := json.Marshal(message)
	if err != nil {
		return err
	}

	err = p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Value: messageBytes,
	}, nil)
	if err != nil {
		return err
	}

	p.Flush(15 * 1000)

	return nil
}

func main() {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost:9092"})
	if err != nil {
		log.Fatal()
	}

	defer p.Close()

	// Delivery report handler for produced messages
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", ev.TopicPartition)
				} else {
					fmt.Printf("Delivered message to %v\n", ev.TopicPartition)
				}
			}
		}
	}()

	messagesJSON, err := readJSONFile("data/messages.json")
	if err != nil {
		log.Fatal(err)
	}

	transactionsJSON, err := readJSONFile("data/transactions.json")
	if err != nil {
		log.Fatal(err)
	}

	var messages []model.Message
	if err := json.Unmarshal(messagesJSON, &messages); err != nil {
		log.Fatal(err)
	}

	var transactions []model.Transaction
	if err := json.Unmarshal(transactionsJSON, &transactions); err != nil {
		log.Fatal(err)
	}

	for _, message := range messages {
		if err := produceMessage(p, "messages", message); err != nil {
			log.Fatal(err)
		}
	}

	for _, transaction := range transactions {
		if err := produceMessage(p, "transactions", transaction); err != nil {
			log.Fatal(err)
		}
	}
}
