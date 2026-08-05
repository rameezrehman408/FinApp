package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"

	"github.com/seisram/FinApp/internal/domain/account"
	"github.com/seisram/FinApp/internal/domain/category"
	"github.com/seisram/FinApp/internal/domain/party"
	"github.com/seisram/FinApp/internal/domain/transaction"
)

type AccountRepository struct {
	db *pgxpool.Pool
}

func NewAccountRepository(db *pgxpool.Pool) *AccountRepository {
	return &AccountRepository{db: db}
}

func (r *AccountRepository) Create(ctx context.Context, a *account.Account) error {
	query := `
		INSERT INTO accounts (id, name, type, currency, balance, description, created_at, updated_at, version)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err := r.db.Exec(ctx, query,
		a.ID, a.Name, a.Type, a.Currency, a.Balance, a.Description,
		a.CreatedAt, a.UpdatedAt, a.Version,
	)
	return err
}

func (r *AccountRepository) GetByID(ctx context.Context, id uuid.UUID) (*account.Account, error) {
	query := `
		SELECT id, name, type, currency, balance, description, created_at, updated_at, version
		FROM accounts WHERE id = $1
	`
	row := r.db.QueryRow(ctx, query, id)
	return r.scanAccount(row)
}

func (r *AccountRepository) GetAll(ctx context.Context) ([]*account.Account, error) {
	query := `
		SELECT id, name, type, currency, balance, description, created_at, updated_at, version
		FROM accounts ORDER BY created_at DESC
	`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []*account.Account
	for rows.Next() {
		a, err := r.scanAccount(rows)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, a)
	}
	return accounts, rows.Err()
}

func (r *AccountRepository) Update(ctx context.Context, a *account.Account) error {
	query := `
		UPDATE accounts SET name = $2, type = $3, currency = $4, balance = $5, 
			description = $6, updated_at = $7, version = $8
		WHERE id = $1 AND version = $8 - 1
	`
	result, err := r.db.Exec(ctx, query,
		a.ID, a.Name, a.Type, a.Currency, a.Balance, a.Description,
		a.UpdatedAt, a.Version,
	)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("account not found or version conflict")
	}
	return nil
}

func (r *AccountRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM accounts WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *AccountRepository) scanAccount(scanner interface{ Scan(dest ...any) error }) (*account.Account, error) {
	var a account.Account
	var balanceStr string
	err := scanner.Scan(
		&a.ID, &a.Name, &a.Type, &a.Currency, &balanceStr,
		&a.Description, &a.CreatedAt, &a.UpdatedAt, &a.Version,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	a.Balance, _ = decimal.NewFromString(balanceStr)
	return &a, nil
}

type TransactionRepository struct {
	db *pgxpool.Pool
}

func NewTransactionRepository(db *pgxpool.Pool) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) Create(ctx context.Context, t *transaction.Transaction) error {
	query := `
		INSERT INTO transactions (id, account_id, counter_account_id, category_id, party_id, type, status, amount, currency, description, reference, transaction_date, value_date, created_at, updated_at, version, metadata)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)
	`
	_, err := r.db.Exec(ctx, query,
		t.ID, t.AccountID, t.CounterAccountID, t.CategoryID, t.PartyID,
		t.Type, t.Status, t.Amount, t.Currency, t.Description, t.Reference,
		t.TransactionDate, t.ValueDate, t.CreatedAt, t.UpdatedAt, t.Version,
		t.Metadata,
	)
	return err
}

func (r *TransactionRepository) GetByID(ctx context.Context, id uuid.UUID) (*transaction.Transaction, error) {
	query := `
		SELECT id, account_id, counter_account_id, category_id, party_id, type, status, amount, currency, description, reference, transaction_date, value_date, created_at, updated_at, version, metadata
		FROM transactions WHERE id = $1
	`
	row := r.db.QueryRow(ctx, query, id)
	return r.scanTransaction(row)
}

func (r *TransactionRepository) GetByAccountID(ctx context.Context, accountID uuid.UUID, limit, offset int) ([]*transaction.Transaction, error) {
	query := `
		SELECT id, account_id, counter_account_id, category_id, party_id, type, status, amount, currency, description, reference, transaction_date, value_date, created_at, updated_at, version, metadata
		FROM transactions WHERE account_id = $1 ORDER BY transaction_date DESC LIMIT $2 OFFSET $3
	`
	rows, err := r.db.Query(ctx, query, accountID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []*transaction.Transaction
	for rows.Next() {
		t, err := r.scanTransaction(rows)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, t)
	}
	return transactions, rows.Err()
}

func (r *TransactionRepository) Update(ctx context.Context, t *transaction.Transaction) error {
	query := `
		UPDATE transactions SET counter_account_id = $2, category_id = $3, party_id = $4, type = $5, status = $6, amount = $7, currency = $8, description = $9, reference = $10, transaction_date = $11, value_date = $12, updated_at = $13, version = $14, metadata = $15
		WHERE id = $1 AND version = $14 - 1
	`
	result, err := r.db.Exec(ctx, query,
		t.ID, t.CounterAccountID, t.CategoryID, t.PartyID, t.Type, t.Status, t.Amount, t.Currency,
		t.Description, t.Reference, t.TransactionDate, t.ValueDate, t.UpdatedAt, t.Version, t.Metadata,
	)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("transaction not found or version conflict")
	}
	return nil
}

func (r *TransactionRepository) scanTransaction(scanner interface{ Scan(dest ...any) error }) (*transaction.Transaction, error) {
	var t transaction.Transaction
	var amountStr string
	var counterAccountID, categoryID, partyID sql.NullString
	var reference sql.NullString
	var metadata map[string]string

	err := scanner.Scan(
		&t.ID, &t.AccountID, &counterAccountID, &categoryID, &partyID,
		&t.Type, &t.Status, &amountStr, &t.Currency, &t.Description, &reference,
		&t.TransactionDate, &t.ValueDate, &t.CreatedAt, &t.UpdatedAt, &t.Version, &metadata,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	t.Amount, _ = decimal.NewFromString(amountStr)
	if counterAccountID.Valid {
		id, _ := uuid.Parse(counterAccountID.String)
		t.CounterAccountID = &id
	}
	if categoryID.Valid {
		id, _ := uuid.Parse(categoryID.String)
		t.CategoryID = &id
	}
	if partyID.Valid {
		id, _ := uuid.Parse(partyID.String)
		t.PartyID = &id
	}
	if reference.Valid {
		t.Reference = reference.String
	}
	t.Metadata = metadata
	return &t, nil
}

type CategoryRepository struct {
	db *pgxpool.Pool
}

func NewCategoryRepository(db *pgxpool.Pool) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) Create(ctx context.Context, c *category.Category) error {
	query := `
		INSERT INTO categories (id, name, type, parent_id, description, icon, color, is_system, created_at, updated_at, version)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`
	_, err := r.db.Exec(ctx, query,
		c.ID, c.Name, c.Type, c.ParentID, c.Description, c.Icon, c.Color, c.IsSystem,
		c.CreatedAt, c.UpdatedAt, c.Version,
	)
	return err
}

func (r *CategoryRepository) GetByID(ctx context.Context, id uuid.UUID) (*category.Category, error) {
	query := `SELECT id, name, type, parent_id, description, icon, color, is_system, created_at, updated_at, version FROM categories WHERE id = $1`
	row := r.db.QueryRow(ctx, query, id)
	return r.scanCategory(row)
}

func (r *CategoryRepository) GetAll(ctx context.Context) ([]*category.Category, error) {
	query := `SELECT id, name, type, parent_id, description, icon, color, is_system, created_at, updated_at, version FROM categories ORDER BY name`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []*category.Category
	for rows.Next() {
		c, err := r.scanCategory(rows)
		if err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	return categories, rows.Err()
}

func (r *CategoryRepository) Update(ctx context.Context, c *category.Category) error {
	query := `
		UPDATE categories SET name = $2, type = $3, parent_id = $4, description = $5, icon = $6, color = $7, updated_at = $8, version = $9
		WHERE id = $1 AND version = $9 - 1
	`
	result, err := r.db.Exec(ctx, query,
		c.ID, c.Name, c.Type, c.ParentID, c.Description, c.Icon, c.Color, c.UpdatedAt, c.Version,
	)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("category not found or version conflict")
	}
	return nil
}

func (r *CategoryRepository) scanCategory(scanner interface{ Scan(dest ...any) error }) (*category.Category, error) {
	var c category.Category
	var parentID sql.NullString
	err := scanner.Scan(
		&c.ID, &c.Name, &c.Type, &parentID, &c.Description, &c.Icon, &c.Color, &c.IsSystem,
		&c.CreatedAt, &c.UpdatedAt, &c.Version,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	if parentID.Valid {
		id, _ := uuid.Parse(parentID.String)
		c.ParentID = &id
	}
	return &c, nil
}

type PartyRepository struct {
	db *pgxpool.Pool
}

func NewPartyRepository(db *pgxpool.Pool) *PartyRepository {
	return &PartyRepository{db: db}
}

func (r *PartyRepository) Create(ctx context.Context, p *party.Party) error {
	query := `
		INSERT INTO parties (id, name, type, tax_id, email, phone, street, city, state, postal_code, country, contact_person, notes, is_active, created_at, updated_at, version)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)
	`
	_, err := r.db.Exec(ctx, query,
		p.ID, p.Name, p.Type, p.TaxID, p.Email, p.Phone,
		p.Address.Street, p.Address.City, p.Address.State, p.Address.PostalCode, p.Address.Country,
		p.ContactPerson, p.Notes, p.IsActive, p.CreatedAt, p.UpdatedAt, p.Version,
	)
	return err
}

func (r *PartyRepository) GetByID(ctx context.Context, id uuid.UUID) (*party.Party, error) {
	query := `
		SELECT id, name, type, tax_id, email, phone, street, city, state, postal_code, country, contact_person, notes, is_active, created_at, updated_at, version
		FROM parties WHERE id = $1
	`
	row := r.db.QueryRow(ctx, query, id)
	return r.scanParty(row)
}

func (r *PartyRepository) GetAll(ctx context.Context) ([]*party.Party, error) {
	query := `SELECT id, name, type, tax_id, email, phone, street, city, state, postal_code, country, contact_person, notes, is_active, created_at, updated_at, version FROM parties ORDER BY name`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var parties []*party.Party
	for rows.Next() {
		p, err := r.scanParty(rows)
		if err != nil {
			return nil, err
		}
		parties = append(parties, p)
	}
	return parties, rows.Err()
}

func (r *PartyRepository) Update(ctx context.Context, p *party.Party) error {
	query := `
		UPDATE parties SET name = $2, type = $3, tax_id = $4, email = $5, phone = $6,
			street = $7, city = $8, state = $9, postal_code = $10, country = $11,
			contact_person = $12, notes = $13, is_active = $14, updated_at = $15, version = $16
		WHERE id = $1 AND version = $16 - 1
	`
	result, err := r.db.Exec(ctx, query,
		p.ID, p.Name, p.Type, p.TaxID, p.Email, p.Phone,
		p.Address.Street, p.Address.City, p.Address.State, p.Address.PostalCode, p.Address.Country,
		p.ContactPerson, p.Notes, p.IsActive, p.UpdatedAt, p.Version,
	)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("party not found or version conflict")
	}
	return nil
}

func (r *PartyRepository) scanParty(scanner interface{ Scan(dest ...any) error }) (*party.Party, error) {
	var p party.Party
	err := scanner.Scan(
		&p.ID, &p.Name, &p.Type, &p.TaxID, &p.Email, &p.Phone,
		&p.Address.Street, &p.Address.City, &p.Address.State, &p.Address.PostalCode, &p.Address.Country,
		&p.ContactPerson, &p.Notes, &p.IsActive, &p.CreatedAt, &p.UpdatedAt, &p.Version,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &p, nil
}