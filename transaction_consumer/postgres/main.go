package main

import (
	"github.com/mmm-bank/financial-analysis-service/storage"
	handle "github.com/mmm-bank/financial-analysis-service/transaction_consumer"
	"github.com/mmm-bank/infra/rabbitmq"
	"log"
	"os"
	"strconv"
)

func main() {
	workerPoolSize, err := strconv.Atoi(os.Getenv("WORKER_POOL_SIZE"))
	if err != nil {
		log.Fatalf("invalid WORKPOOL_SIZE: %v", err)
	}

	s := storage.NewPostgresTransactions(os.Getenv("POSTGRES_URL"))
	conn := messaging.NewConn(os.Getenv("RABBITMQ_URL"))
	if err := messaging.DeclareQueue("postgres", conn); err != nil {
		log.Fatalf("failed to declare queue: %v", err)
	}
	for i := 0; i < workerPoolSize-1; i++ {
		go handle.Events("postgres", conn, s)
	}
	handle.Events("postgres", conn, s)
}
