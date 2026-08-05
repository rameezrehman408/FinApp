package query

import (
	"context"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/seisram/FinApp/internal/domain/account"
	"github.com/seisram/FinApp/internal/domain/category"
	"github.com/seisram/FinApp/internal/domain/party"
	"github.com/seisram/FinApp/internal/domain/transaction"
)

type AccountQueryHandler struct {
	repo AccountRepository
}

type AccountRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*account.Account, error)
	GetAll(ctx context.Context) ([]*account.Account, error)
	GetByType(ctx context.Context, accountType account.AccountType) ([]*account.Account, error)
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

type TransactionQueryHandler struct {
	repo TransactionRepository
}

type TransactionRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*transaction.Transaction, error)
	GetByAccountID(ctx context.Context, accountID uuid.UUID, limit, offset int) ([]*transaction.Transaction, error)
	GetByDateRange(ctx context.Context, accountID uuid.UUID, startDate, endDate time.Time) ([]*transaction.Transaction, error)
	GetByCategoryID(ctx context.Context, categoryID uuid.UUID, limit, offset int) ([]*transaction.Transaction, error)
	GetByPartyID(ctx context.Context, partyID uuid.UUID, limit, offset int) ([]*transaction.Transaction, error)
	GetByStatus(ctx context.Context, status transaction.TransactionStatus, limit, offset int) ([]*transaction.Transaction, error)
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
	} else if q.CategoryID != nil {
		transactions, err = h.repo.GetByCategoryID(ctx, *q.CategoryID, q.Limit, q.Offset)
	} else if q.PartyID != nil {
		transactions, err = h.repo.GetByPartyID(ctx, *q.PartyID, q.Limit, q.Offset)
	} else if q.Status != nil {
		transactions, err = h.repo.GetByStatus(ctx, *q.Status, q.Limit, q.Offset)
	} else {
		// Default: get all transactions (would need a GetAll method)
		transactions = []*transaction.Transaction{}
	}

	if err != nil {
		return nil, err
	}

	return &ListTransactionsResult{
		Transactions: transactions,
		Total:        len(transactions),
	}, nil
}

type CategoryQueryHandler struct {
	repo CategoryRepository
}

type CategoryRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*category.Category, error)
	GetAll(ctx context.Context) ([]*category.Category, error)
	GetByType(ctx context.Context, categoryType category.CategoryType) ([]*category.Category, error)
	GetChildren(ctx context.Context, parentID uuid.UUID) ([]*category.Category, error)
}

func NewCategoryQueryHandler(repo CategoryRepository) *CategoryQueryHandler {
	return &CategoryQueryHandler{repo: repo}
}

type GetCategoryQuery struct {
	ID uuid.UUID
}

type GetCategoryResult struct {
	Category *category.Category
}

func (h *CategoryQueryHandler) HandleGetCategory(ctx context.Context, q GetCategoryQuery) (*GetCategoryResult, error) {
	c, err := h.repo.GetByID(ctx, q.ID)
	if err != nil {
		return nil, err
	}
	return &GetCategoryResult{Category: c}, nil
}

type ListCategoriesQuery struct {
	Type     *category.CategoryType
	ParentID *uuid.UUID
}

type ListCategoriesResult struct {
	Categories []*category.Category
}

func (h *CategoryQueryHandler) HandleListCategories(ctx context.Context, q ListCategoriesQuery) (*ListCategoriesResult, error) {
	var categories []*category.Category
	var err error

	if q.ParentID != nil {
		categories, err = h.repo.GetChildren(ctx, *q.ParentID)
	} else if q.Type != nil {
		categories, err = h.repo.GetByType(ctx, *q.Type)
	} else {
		categories, err = h.repo.GetAll(ctx)
	}

	if err != nil {
		return nil, err
	}
	return &ListCategoriesResult{Categories: categories}, nil
}

type PartyQueryHandler struct {
	repo PartyRepository
}

type PartyRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*party.Party, error)
	GetAll(ctx context.Context) ([]*party.Party, error)
	GetByType(ctx context.Context, partyType party.PartyType) ([]*party.Party, error)
	GetReceivables(ctx context.Context, partyID uuid.UUID) ([]*party.Receivable, error)
	GetPayables(ctx context.Context, partyID uuid.UUID) ([]*party.Payable, error)
	GetOverdueReceivables(ctx context.Context) ([]*party.Receivable, error)
	GetOverduePayables(ctx context.Context) ([]*party.Payable, error)
}

func NewPartyQueryHandler(repo PartyRepository) *PartyQueryHandler {
	return &PartyQueryHandler{repo: repo}
}

type GetPartyQuery struct {
	ID uuid.UUID
}

type GetPartyResult struct {
	Party *party.Party
}

func (h *PartyQueryHandler) HandleGetParty(ctx context.Context, q GetPartyQuery) (*GetPartyResult, error) {
	p, err := h.repo.GetByID(ctx, q.ID)
	if err != nil {
		return nil, err
	}
	return &GetPartyResult{Party: p}, nil
}

type ListPartiesQuery struct {
	Type *party.PartyType
}

type ListPartiesResult struct {
	Parties []*party.Party
}

func (h *PartyQueryHandler) HandleListParties(ctx context.Context, q ListPartiesQuery) (*ListPartiesResult, error) {
	var parties []*party.Party
	var err error

	if q.Type != nil {
		parties, err = h.repo.GetByType(ctx, *q.Type)
	} else {
		parties, err = h.repo.GetAll(ctx)
	}

	if err != nil {
		return nil, err
	}
	return &ListPartiesResult{Parties: parties}, nil
}

type GetReceivablesQuery struct {
	PartyID uuid.UUID
}

type GetReceivablesResult struct {
	Receivables []*party.Receivable
}

func (h *PartyQueryHandler) HandleGetReceivables(ctx context.Context, q GetReceivablesQuery) (*GetReceivablesResult, error) {
	r, err := h.repo.GetReceivables(ctx, q.PartyID)
	if err != nil {
		return nil, err
	}
	return &GetReceivablesResult{Receivables: r}, nil
}

type GetPayablesQuery struct {
	PartyID uuid.UUID
}

type GetPayablesResult struct {
	Payables []*party.Payable
}

func (h *PartyQueryHandler) HandleGetPayables(ctx context.Context, q GetPayablesQuery) (*GetPayablesResult, error) {
	p, err := h.repo.GetPayables(ctx, q.PartyID)
	if err != nil {
		return nil, err
	}
	return &GetPayablesResult{Payables: p}, nil
}

type GetOverdueReceivablesQuery struct{}

type GetOverdueReceivablesResult struct {
	Receivables []*party.Receivable
}

func (h *PartyQueryHandler) HandleGetOverdueReceivables(ctx context.Context, q GetOverdueReceivablesQuery) (*GetOverdueReceivablesResult, error) {
	r, err := h.repo.GetOverdueReceivables(ctx)
	if err != nil {
		return nil, err
	}
	return &GetOverdueReceivablesResult{Receivables: r}, nil
}

type GetOverduePayablesQuery struct{}

type GetOverduePayablesResult struct {
	Payables []*party.Payable
}

func (h *PartyQueryHandler) HandleGetOverduePayables(ctx context.Context, q GetOverduePayablesQuery) (*GetOverduePayablesResult, error) {
	p, err := h.repo.GetOverduePayables(ctx)
	if err != nil {
		return nil, err
	}
	return &GetOverduePayablesResult{Payables: p}, nil
}