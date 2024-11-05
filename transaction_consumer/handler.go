package handle

import (
	"encoding/json"
	"fmt"
	"github.com/mmm-bank/financial-analysis-service/models"
	messaging "github.com/mmm-bank/infra/rabbitmq"
	"github.com/streadway/amqp"
)

type Storage interface {
	AddTransfer(transfer *models.Transfer) error
}

func Events(queueName string, conn *amqp.Connection, storage Storage) {
	eventHandler := messaging.NewConsumer(conn)
	err := eventHandler.ConsumeMessages(queueName, func(message []byte) error {
		transaction := &models.Transfer{}
		if err := json.Unmarshal(message, transaction); err != nil {
			return fmt.Errorf("error decoding JSON: %v", err)
		}
		if err := storage.AddTransfer(transaction); err != nil {
			return fmt.Errorf("failed to add transaction: %v", err)
		}
		return nil
	})
	if err != nil {
		fmt.Printf("error consuming transaction messages: %s\n", err)
	}
}
