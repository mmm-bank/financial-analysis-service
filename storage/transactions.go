package storage

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mmm-bank/financial-analysis-service/models"
	"log"
	"time"
)

var _ TransactionStorage = PostgresTransactions{}

type TransactionStorage interface {
	GetTransactions(userID uuid.UUID) ([]models.TransactionInfo, error)
}

type PostgresTransactions struct {
	db *pgxpool.Pool
}

func NewPostgresTransactions(connString string) PostgresTransactions {
	pool, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v", err)
	}
	return PostgresTransactions{pool}
}

func (p PostgresTransactions) GetTransactions(userID uuid.UUID) ([]models.TransactionInfo, error) {
	query := "SELECT account_id, partner_account_id, transaction_type, amount FROM transactions WHERE user_id = $1 ORDER BY created_at DESC"
	rows, err := p.db.Query(context.Background(), query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get all transactions: %v", err)
	}
	defer rows.Close()

	var transactions []models.TransactionInfo
	for rows.Next() {
		var transaction models.TransactionInfo
		err := rows.Scan(&transaction.AccountID, &transaction.PartnerAccountID, &transaction.TransactionType, &transaction.Amount)
		if err != nil {
			return nil, fmt.Errorf("failed to scan transaction: %v", err)
		}
		transactions = append(transactions, transaction)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("error during rows iteration: %v", rows.Err())
	}

	return transactions, nil
}

func (t PostgresTransactions) AddTransfer(transfer *models.Transfer) error {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	_, err := t.db.Exec(ctx, "CALL add_transfer($1, $2, $3, $4, $5, $6)",
		transfer.ID,
		transfer.SenderID,
		transfer.SenderAccountID,
		transfer.ReceiverID,
		transfer.ReceiverAccountID,
		transfer.Amount)
	return err
}
