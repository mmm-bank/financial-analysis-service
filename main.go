package main

import (
	. "github.com/mmm-bank/financial-analysis-service/http"
	. "github.com/mmm-bank/financial-analysis-service/storage"
	"log"
	"os"
)

func main() {
	addr := ":8080"
	s := NewPostgresTransactions(os.Getenv("POSTGRES_URL"))
	m := NewMongoAnalysis(os.Getenv("MONGO_URL"))
	server := NewTransactionService(s, m)

	log.Printf("Expense analysis server is running on port %s...", addr[1:])
	if err := CreateAndRunServer(server, addr); err != nil {
		log.Fatalf("Failed to start server: %s", err)
	}
}
