package transaction

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type TransactionType string

const (
	TransactionTypeIncome     TransactionType = "INCOME"
	TransactionTypeExpense    TransactionType = "EXPENSE"
	TransactionTypeTransfer   TransactionType = "TRANSFER"
	TransactionTypeInvestment TransactionType = "INVESTMENT"
	TransactionTypeReceivable TransactionType = "RECEIVABLE"
	TransactionTypePayable    TransactionType = "PAYABLE"
)

type TransactionStatus string

const (
	TransactionStatusPending   TransactionStatus = "PENDING"
	TransactionStatusCompleted TransactionStatus = "COMPLETED"
	TransactionStatusCancelled TransactionStatus = "CANCELLED"
	TransactionStatusReversed  TransactionStatus = "REVERSED"
)

type Transaction struct {
	ID              uuid.UUID
	AccountID       uuid.UUID
	CounterAccountID *uuid.UUID // For transfers
	CategoryID      *uuid.UUID
	PartyID         *uuid.UUID
	Type            TransactionType
	Status          TransactionStatus
	Amount          decimal.Decimal
	Currency        string
	Description     string
	Reference       string // External reference (invoice #, transaction ID, etc.)
	TransactionDate time.Time
	ValueDate       time.Time // When the transaction actually affects balance
	CreatedAt       time.Time
	UpdatedAt       time.Time
	Version         int
	Metadata        map[string]string // Flexible metadata for extensibility
}

func NewTransaction(
	accountID uuid.UUID,
	transactionType TransactionType,
	amount decimal.Decimal,
	currency string,
	description string,
	transactionDate time.Time,
) *Transaction {
	now := time.Now().UTC()
	return &Transaction{
		ID:              uuid.New(),
		AccountID:       accountID,
		Type:            transactionType,
		Status:          TransactionStatusPending,
		Amount:          amount,
		Currency:        currency,
		Description:     description,
		TransactionDate: transactionDate,
		ValueDate:       transactionDate,
		CreatedAt:       now,
		UpdatedAt:       now,
		Version:         1,
		Metadata:        make(map[string]string),
	}
}

func (t *Transaction) Complete() error {
	if t.Status != TransactionStatusPending {
		return ErrInvalidStatusTransition
	}
	t.Status = TransactionStatusCompleted
	t.UpdatedAt = time.Now().UTC()
	t.Version++
	return nil
}

func (t *Transaction) Cancel() error {
	if t.Status == TransactionStatusCompleted {
		return ErrCannotCancelCompleted
	}
	if t.Status == TransactionStatusCancelled {
		return ErrAlreadyCancelled
	}
	t.Status = TransactionStatusCancelled
	t.UpdatedAt = time.Now().UTC()
	t.Version++
	return nil
}

func (t *Transaction) Reverse() *Transaction {
	reversal := NewTransaction(
		t.AccountID,
		t.Type,
		t.Amount.Neg(),
		t.Currency,
		"Reversal: "+t.Description,
		time.Now().UTC(),
	)
	reversal.CounterAccountID = t.CounterAccountID
	reversal.CategoryID = t.CategoryID
	reversal.PartyID = t.PartyID
	reversal.Reference = "REV-" + t.Reference
	reversal.Metadata["original_transaction_id"] = t.ID.String()
	reversal.Metadata["reversal_reason"] = "MANUAL_REVERSAL"
	return reversal
}

var (
	ErrInvalidStatusTransition = &DomainError{Code: "INVALID_STATUS_TRANSITION", Message: "Invalid status transition"}
	ErrCannotCancelCompleted   = &DomainError{Code: "CANNOT_CANCEL_COMPLETED", Message: "Cannot cancel a completed transaction"}
	ErrAlreadyCancelled        = &DomainError{Code: "ALREADY_CANCELLED", Message: "Transaction already cancelled"}
)

type DomainError struct {
	Code    string
	Message string
}

func (e *DomainError) Error() string {
	return e.Message
}