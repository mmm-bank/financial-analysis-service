package models

import "github.com/google/uuid"

type Transfer struct {
	ID                uuid.UUID `json:"transaction_id"`
	SenderID          uuid.UUID `json:"sender_id"`
	SenderAccountID   uuid.UUID `json:"sender_account_id"`
	ReceiverID        uuid.UUID `json:"receiver_id"`
	ReceiverAccountID uuid.UUID `json:"receiver_account_id"`
	Amount            uint64    `json:"amount"`
}
