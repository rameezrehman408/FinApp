package account

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type AccountType string

const (
	AccountTypeAsset     AccountType = "ASSET"
	AccountTypeLiability AccountType = "LIABILITY"
	AccountTypeEquity    AccountType = "EQUITY"
	AccountTypeIncome    AccountType = "INCOME"
	AccountTypeExpense   AccountType = "EXPENSE"
)

type Account struct {
	ID          uuid.UUID
	Name        string
	Type        AccountType
	Currency    string
	Balance     decimal.Decimal
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Version     int
}

func NewAccount(name string, accountType AccountType, currency string, description string) *Account {
	now := time.Now().UTC()
	return &Account{
		ID:          uuid.New(),
		Name:        name,
		Type:        accountType,
		Currency:    currency,
		Balance:     decimal.Zero,
		Description: description,
		CreatedAt:   now,
		UpdatedAt:   now,
		Version:     1,
	}
}

func (a *Account) Credit(amount decimal.Decimal) error {
	if amount.LessThanOrEqual(decimal.Zero) {
		return ErrInvalidAmount
	}
	a.Balance = a.Balance.Add(amount)
	a.UpdatedAt = time.Now().UTC()
	a.Version++
	return nil
}

func (a *Account) Debit(amount decimal.Decimal) error {
	if amount.LessThanOrEqual(decimal.Zero) {
		return ErrInvalidAmount
	}
	// For liability/equity/income accounts, balance can go negative
	// For asset/expense accounts, we might want to prevent negative balance
	// This is a business rule decision - for now allow negative
	a.Balance = a.Balance.Sub(amount)
	a.UpdatedAt = time.Now().UTC()
	a.Version++
	return nil
}

var ErrInvalidAmount = &DomainError{Code: "INVALID_AMOUNT", Message: "Amount must be positive"}

type DomainError struct {
	Code    string
	Message string
}

func (e *DomainError) Error() string {
	return e.Message
}