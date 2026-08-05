package command

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/seisram/FinApp/internal/domain/account"
	"github.com/seisram/FinApp/internal/domain/category"
	"github.com/seisram/FinApp/internal/domain/party"
	"github.com/seisram/FinApp/internal/domain/transaction"
)

type AccountRepository interface {
	Create(ctx context.Context, a *account.Account) error
	GetByID(ctx context.Context, id uuid.UUID) (*account.Account, error)
	Update(ctx context.Context, a *account.Account) error
}

type CreateAccountHandler struct {
	repo AccountRepository
}

func NewCreateAccountHandler(repo AccountRepository) *CreateAccountHandler {
	return &CreateAccountHandler{repo: repo}
}

type CreateAccountCommand struct {
	Name        string
	Type        account.AccountType
	Currency    string
	Description string
}

type CreateAccountResult struct {
	Account *account.Account
}

func (h *CreateAccountHandler) Handle(ctx context.Context, cmd CreateAccountCommand) (*CreateAccountResult, error) {
	a := account.NewAccount(cmd.Name, cmd.Type, cmd.Currency, cmd.Description)
	if err := h.repo.Create(ctx, a); err != nil {
		return nil, err
	}
	return &CreateAccountResult{Account: a}, nil
}

type TransactionRepository interface {
	Create(ctx context.Context, t *transaction.Transaction) error
	GetByID(ctx context.Context, id uuid.UUID) (*transaction.Transaction, error)
	Update(ctx context.Context, t *transaction.Transaction) error
}

type AccountRepositoryForTransaction interface {
	GetByID(ctx context.Context, id uuid.UUID) (*account.Account, error)
	Update(ctx context.Context, a *account.Account) error
}

type CreateTransactionHandler struct {
	repo         TransactionRepository
	accountRepo  AccountRepositoryForTransaction
}

func NewCreateTransactionHandler(repo TransactionRepository, accountRepo AccountRepositoryForTransaction) *CreateTransactionHandler {
	return &CreateTransactionHandler{repo: repo, accountRepo: accountRepo}
}

type CreateTransactionCommand struct {
	AccountID       uuid.UUID
	CounterAccountID *uuid.UUID
	CategoryID      *uuid.UUID
	PartyID         *uuid.UUID
	Type            transaction.TransactionType
	Amount          decimal.Decimal
	Currency        string
	Description     string
	Reference       string
	TransactionDate time.Time
}

type CreateTransactionResult struct {
	Transaction *transaction.Transaction
}

func (h *CreateTransactionHandler) Handle(ctx context.Context, cmd CreateTransactionCommand) (*CreateTransactionResult, error) {
	// Validate account exists
	acc, err := h.accountRepo.GetByID(ctx, cmd.AccountID)
	if err != nil {
		return nil, err
	}
	if acc == nil {
		return nil, ErrAccountNotFound
	}

	t := transaction.NewTransaction(
		cmd.AccountID,
		cmd.Type,
		cmd.Amount,
		cmd.Currency,
		cmd.Description,
		cmd.TransactionDate,
	)
	t.CounterAccountID = cmd.CounterAccountID
	t.CategoryID = cmd.CategoryID
	t.PartyID = cmd.PartyID
	t.Reference = cmd.Reference

	if err := h.repo.Create(ctx, t); err != nil {
		return nil, err
	}

	// Update account balance based on transaction type
	switch cmd.Type {
	case transaction.TransactionTypeIncome:
		if err := acc.Credit(cmd.Amount); err != nil {
			return nil, err
		}
	case transaction.TransactionTypeExpense:
		if err := acc.Debit(cmd.Amount); err != nil {
			return nil, err
		}
	case transaction.TransactionTypeTransfer:
		if cmd.CounterAccountID != nil {
			counterAcc, err := h.accountRepo.GetByID(ctx, *cmd.CounterAccountID)
			if err != nil {
				return nil, err
			}
			if counterAcc == nil {
				return nil, ErrCounterAccountNotFound
			}
			if err := acc.Debit(cmd.Amount); err != nil {
				return nil, err
			}
			if err := counterAcc.Credit(cmd.Amount); err != nil {
				return nil, err
			}
			if err := h.accountRepo.Update(ctx, counterAcc); err != nil {
				return nil, err
			}
		}
	}

	if err := h.accountRepo.Update(ctx, acc); err != nil {
		return nil, err
	}

	return &CreateTransactionResult{Transaction: t}, nil
}

var (
	ErrAccountNotFound       = &CommandError{Code: "ACCOUNT_NOT_FOUND", Message: "Account not found"}
	ErrCounterAccountNotFound = &CommandError{Code: "COUNTER_ACCOUNT_NOT_FOUND", Message: "Counter account not found"}
)

type CommandError struct {
	Code    string
	Message string
}

func (e *CommandError) Error() string {
	return e.Message
}

type CategoryRepository interface {
	Create(ctx context.Context, c *category.Category) error
	GetByID(ctx context.Context, id uuid.UUID) (*category.Category, error)
	Update(ctx context.Context, c *category.Category) error
}

type CreateCategoryHandler struct {
	repo CategoryRepository
}

func NewCreateCategoryHandler(repo CategoryRepository) *CreateCategoryHandler {
	return &CreateCategoryHandler{repo: repo}
}

type CreateCategoryCommand struct {
	Name        string
	Type        category.CategoryType
	ParentID    *uuid.UUID
	Description string
	Icon        string
	Color       string
}

type CreateCategoryResult struct {
	Category *category.Category
}

func (h *CreateCategoryHandler) Handle(ctx context.Context, cmd CreateCategoryCommand) (*CreateCategoryResult, error) {
	c := category.NewCategory(cmd.Name, cmd.Type, cmd.ParentID, cmd.Description)
	c.Icon = cmd.Icon
	c.Color = cmd.Color
	if err := h.repo.Create(ctx, c); err != nil {
		return nil, err
	}
	return &CreateCategoryResult{Category: c}, nil
}

type PartyRepository interface {
	Create(ctx context.Context, p *party.Party) error
	GetByID(ctx context.Context, id uuid.UUID) (*party.Party, error)
	Update(ctx context.Context, p *party.Party) error
}

type CreatePartyHandler struct {
	repo PartyRepository
}

func NewCreatePartyHandler(repo PartyRepository) *CreatePartyHandler {
	return &CreatePartyHandler{repo: repo}
}

type CreatePartyCommand struct {
	Name         string
	Type         party.PartyType
	TaxID        string
	Email        string
	Phone        string
	Street       string
	City         string
	State        string
	PostalCode   string
	Country      string
	ContactPerson string
	Notes        string
}

type CreatePartyResult struct {
	Party *party.Party
}

func (h *CreatePartyHandler) Handle(ctx context.Context, cmd CreatePartyCommand) (*CreatePartyResult, error) {
	p := party.NewParty(cmd.Name, cmd.Type)
	p.TaxID = cmd.TaxID
	p.Email = cmd.Email
	p.Phone = cmd.Phone
	p.Address.Street = cmd.Street
	p.Address.City = cmd.City
	p.Address.State = cmd.State
	p.Address.PostalCode = cmd.PostalCode
	p.Address.Country = cmd.Country
	p.ContactPerson = cmd.ContactPerson
	p.Notes = cmd.Notes

	if err := h.repo.Create(ctx, p); err != nil {
		return nil, err
	}
	return &CreatePartyResult{Party: p}, nil
}