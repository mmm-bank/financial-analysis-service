package models

import (
	"github.com/google/uuid"
	"time"
)

type ExpensesTransactionInfo struct {
	ID         uuid.UUID `bson:"id"`
	ReceiverID uuid.UUID `bson:"receiver_id"`
	Category   string    `bson:"category"`
	Cost       uint64    `bson:"cost"`
	Timestamp  time.Time `bson:"created_at"`
}

type ExpensesAnalysis struct {
	UserID           uuid.UUID                 `bson:"user_id"`
	MonthYear        string                    `bson:"month_year"`
	TotalCost        uint64                    `bson:"total_cost"`
	TransactionCount uint64                    `bson:"transaction_count"`
	Transactions     []ExpensesTransactionInfo `bson:"transactions"`
}
