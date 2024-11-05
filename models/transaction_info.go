package models

import "github.com/google/uuid"

type TransactionInfo struct {
	AccountID        uuid.UUID `json:"account_id"`
	PartnerAccountID uuid.UUID `json:"partner_account_id"`
	TransactionType  string    `json:"transaction_type"`
	Amount           uint64    `json:"amount"`
}
