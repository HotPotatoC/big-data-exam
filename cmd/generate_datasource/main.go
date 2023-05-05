package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/HotPotatoC/bigdata-exam/model"
	"github.com/google/uuid"
	"github.com/jaswdr/faker"
)

type input struct {
	id        string
	timestamp time.Time
}

func generateMessageData(data []input) []byte {
	fake := faker.New()

	var messages []model.Message

	for i := 0; i < len(data); i++ {
		message := model.NewMessage(
			data[i].id,
			"Hello there i am "+fake.App().Name(),
			data[i].timestamp)
		messages = append(messages, message)
	}

	jsonData, _ := json.Marshal(messages)

	return jsonData
}

func generateTransactionData(data []input) []byte {
	rand.Seed(time.Now().UnixNano())

	var transactions []model.Transaction

	for i := 0; i < len(data); i++ {
		// to make sure that the transaction is not heading to the same user
		transactionToIDX := rand.Intn(len(data))
		for transactionToIDX == i {
			transactionToIDX = rand.Intn(len(data))
		}

		transaction := model.NewTransaction(
			data[i].id,
			data[i].id,
			data[transactionToIDX].id,
			rand.Intn(1000000),
			data[i].timestamp)

		transactions = append(transactions, transaction)
	}

	jsonData, _ := json.Marshal(transactions)

	return jsonData
}

func main() {
	n := 50
	fake := faker.New()

	var data []input

	// Make sure the ID and Timestamp is in sync
	for i := 0; i < n; i++ {
		data = append(data, input{
			id:        uuid.NewString(),
			timestamp: fake.Time().TimeBetween(time.Now().Add(-time.Hour*24*30), time.Now())})
	}

	messagesJSON := generateMessageData(data)
	transactionsJSON := generateTransactionData(data)

	err := os.WriteFile("data/messages.json", messagesJSON, 0777)
	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile("data/transactions.json", transactionsJSON, 0777)
	if err != nil {
		log.Fatal(err)
	}
}
