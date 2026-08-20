package query

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/seisram/FinApp/internal/domain/account"
	"github.com/seisram/FinApp/internal/domain/category"
	"github.com/seisram/FinApp/internal/domain/transaction"
)

type AccountRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*account.Account, error)
	GetAll(ctx context.Context) ([]*account.Account, error)
	GetByType(ctx context.Context, accountType account.AccountType) ([]*account.Account, error)
}

type AccountQueryHandler struct {
	repo AccountRepository
}

func NewAccountQueryHandler(repo AccountRepository) *AccountQueryHandler {
	return &AccountQueryHandler{repo: repo}
}

type GetAccountQuery struct {
	ID uuid.UUID
}

type GetAccountResult struct {
	Account *account.Account
}

func (h *AccountQueryHandler) HandleGetAccount(ctx context.Context, q GetAccountQuery) (*GetAccountResult, error) {
	a, err := h.repo.GetByID(ctx, q.ID)
	if err != nil {
		return nil, err
	}
	return &GetAccountResult{Account: a}, nil
}

type ListAccountsQuery struct {
	Type *account.AccountType
}

type ListAccountsResult struct {
	Accounts []*account.Account
}

func (h *AccountQueryHandler) HandleListAccounts(ctx context.Context, q ListAccountsQuery) (*ListAccountsResult, error) {
	var accounts []*account.Account
	var err error

	if q.Type != nil {
		accounts, err = h.repo.GetByType(ctx, *q.Type)
	} else {
		accounts, err = h.repo.GetAll(ctx)
	}

	if err != nil {
		return nil, err
	}
	return &ListAccountsResult{Accounts: accounts}, nil
}

type TransactionRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*transaction.Transaction, error)
	GetByAccountID(ctx context.Context, accountID uuid.UUID, limit, offset int) ([]*transaction.Transaction, error)
	GetByDateRange(ctx context.Context, accountID uuid.UUID, startDate, endDate time.Time) ([]*transaction.Transaction, error)
	GetByCategoryID(ctx context.Context, categoryID uuid.UUID, limit, offset int) ([]*transaction.Transaction, error)
	GetByPartyID(ctx context.Context, partyID uuid.UUID, limit, offset int) ([]*transaction.Transaction, error)
	GetByStatus(ctx context.Context, status transaction.TransactionStatus, limit, offset int) ([]*transaction.Transaction, error)
}

type TransactionQueryHandler struct {
	repo TransactionRepository
}

func NewTransactionQueryHandler(repo TransactionRepository) *TransactionQueryHandler {
	return &TransactionQueryHandler{repo: repo}
	}

type GetTransactionQuery struct {
	ID uuid.UUID
}

type GetTransactionResult struct {
	Transaction *transaction.Transaction
}

func (h *TransactionQueryHandler) HandleGetTransaction(ctx context.Context, q GetTransactionQuery) (*GetTransactionResult, error) {
	t, err := h.repo.GetByID(ctx, q.ID)
	if err != nil {
		return nil, err
	}
	return &GetTransactionResult{Transaction: t}, nil
}

type ListTransactionsQuery struct {
	AccountID  *uuid.UUID
	CategoryID *uuid.UUID
	PartyID    *uuid.UUID
	Status     *transaction.TransactionStatus
	StartDate  *time.Time
	EndDate    *time.Time
	Limit      int
	Offset     int
}

type ListTransactionsResult struct {
	Transactions []*transaction.Transaction
	Total        int
}

func (h *TransactionQueryHandler) HandleListTransactions(ctx context.Context, q ListTransactionsQuery) (*ListTransactionsResult, error) {
	var transactions []*transaction.Transaction
	var err error

	if q.AccountID != nil {
		transactions, err = h.repo.GetByAccountID(ctx, *q.AccountID, q.Limit, q.Offset)
	}
	return &ListTransactionsResult{Transactions: transactions, Total: len(transactions)}, err
	}

type CategoryRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*category.Category, error)
	GetAll(ctx context.Context) ([]*category.Category, error)
}
