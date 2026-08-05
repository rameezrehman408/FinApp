-- Initial schema for FinApp
-- Migration: 000001_initial_schema.up.sql

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Accounts table
CREATE TABLE accounts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    type VARCHAR(50) NOT NULL CHECK (type IN ('ASSET', 'LIABILITY', 'EQUITY', 'INCOME', 'EXPENSE')),
    currency CHAR(3) NOT NULL DEFAULT 'USD',
    balance NUMERIC(19, 4) NOT NULL DEFAULT 0,
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    version INTEGER NOT NULL DEFAULT 1
);

CREATE INDEX idx_accounts_type ON accounts(type);
CREATE INDEX idx_accounts_currency ON accounts(currency);

-- Categories table
CREATE TABLE categories (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    type VARCHAR(50) NOT NULL CHECK (type IN ('INCOME', 'EXPENSE', 'ASSET', 'LIABILITY')),
    parent_id UUID REFERENCES categories(id) ON DELETE SET NULL,
    description TEXT,
    icon VARCHAR(100),
    color VARCHAR(7),
    is_system BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    version INTEGER NOT NULL DEFAULT 1
);

CREATE INDEX idx_categories_type ON categories(type);
CREATE INDEX idx_categories_parent ON categories(parent_id);

-- Parties table
CREATE TABLE parties (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    type VARCHAR(50) NOT NULL CHECK (type IN ('PERSON', 'COMPANY', 'GOVERNMENT', 'OTHER')),
    tax_id VARCHAR(100),
    email VARCHAR(255),
    phone VARCHAR(50),
    street VARCHAR(255),
    city VARCHAR(100),
    state VARCHAR(100),
    postal_code VARCHAR(20),
    country VARCHAR(100),
    contact_person VARCHAR(255),
    notes TEXT,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    version INTEGER NOT NULL DEFAULT 1
);

CREATE INDEX idx_parties_type ON parties(type);
CREATE INDEX idx_parties_email ON parties(email);
CREATE INDEX idx_parties_is_active ON parties(is_active);

-- Transactions table
CREATE TABLE transactions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    account_id UUID NOT NULL REFERENCES accounts(id) ON DELETE RESTRICT,
    counter_account_id UUID REFERENCES accounts(id) ON DELETE SET NULL,
    category_id UUID REFERENCES categories(id) ON DELETE SET NULL,
    party_id UUID REFERENCES parties(id) ON DELETE SET NULL,
    type VARCHAR(50) NOT NULL CHECK (type IN ('INCOME', 'EXPENSE', 'TRANSFER', 'INVESTMENT', 'RECEIVABLE', 'PAYABLE')),
    status VARCHAR(50) NOT NULL DEFAULT 'PENDING' CHECK (status IN ('PENDING', 'COMPLETED', 'CANCELLED', 'REVERSED')),
    amount NUMERIC(19, 4) NOT NULL,
    currency CHAR(3) NOT NULL DEFAULT 'USD',
    description TEXT,
    reference VARCHAR(255),
    transaction_date TIMESTAMPTZ NOT NULL,
    value_date TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    version INTEGER NOT NULL DEFAULT 1,
    metadata JSONB DEFAULT '{}'
);

CREATE INDEX idx_transactions_account ON transactions(account_id);
CREATE INDEX idx_transactions_counter_account ON transactions(counter_account_id);
CREATE INDEX idx_transactions_category ON transactions(category_id);
CREATE INDEX idx_transactions_party ON transactions(party_id);
CREATE INDEX idx_transactions_type ON transactions(type);
CREATE INDEX idx_transactions_status ON transactions(status);
CREATE INDEX idx_transactions_date ON transactions(transaction_date);
CREATE INDEX idx_transactions_value_date ON transactions(value_date);

-- Receivables table
CREATE TABLE receivables (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    party_id UUID NOT NULL REFERENCES parties(id) ON DELETE RESTRICT,
    account_id UUID NOT NULL REFERENCES accounts(id) ON DELETE RESTRICT,
    invoice_number VARCHAR(100) NOT NULL,
    description TEXT,
    amount NUMERIC(19, 4) NOT NULL,
    currency CHAR(3) NOT NULL DEFAULT 'USD',
    amount_received NUMERIC(19, 4) NOT NULL DEFAULT 0,
    status VARCHAR(50) NOT NULL DEFAULT 'PENDING' CHECK (status IN ('PENDING', 'PARTIAL', 'PAID', 'OVERDUE', 'CANCELLED', 'WRITTEN_OFF')),
    issue_date TIMESTAMPTZ NOT NULL,
    due_date TIMESTAMPTZ NOT NULL,
    paid_date TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    version INTEGER NOT NULL DEFAULT 1,
    metadata JSONB DEFAULT '{}'
);

CREATE INDEX idx_receivables_party ON receivables(party_id);
CREATE INDEX idx_receivables_account ON receivables(account_id);
CREATE INDEX idx_receivables_status ON receivables(status);
CREATE INDEX idx_receivables_due_date ON receivables(due_date);
CREATE UNIQUE INDEX idx_receivables_invoice ON receivables(invoice_number);

-- Payables table
CREATE TABLE payables (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    party_id UUID NOT NULL REFERENCES parties(id) ON DELETE RESTRICT,
    account_id UUID NOT NULL REFERENCES accounts(id) ON DELETE RESTRICT,
    invoice_number VARCHAR(100) NOT NULL,
    description TEXT,
    amount NUMERIC(19, 4) NOT NULL,
    currency CHAR(3) NOT NULL DEFAULT 'USD',
    amount_paid NUMERIC(19, 4) NOT NULL DEFAULT 0,
    status VARCHAR(50) NOT NULL DEFAULT 'PENDING' CHECK (status IN ('PENDING', 'PARTIAL', 'PAID', 'OVERDUE', 'CANCELLED')),
    issue_date TIMESTAMPTZ NOT NULL,
    due_date TIMESTAMPTZ NOT NULL,
    paid_date TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    version INTEGER NOT NULL DEFAULT 1,
    metadata JSONB DEFAULT '{}'
);

CREATE INDEX idx_payables_party ON payables(party_id);
CREATE INDEX idx_payables_account ON payables(account_id);
CREATE INDEX idx_payables_status ON payables(status);
CREATE INDEX idx_payables_due_date ON payables(due_date);
CREATE UNIQUE INDEX idx_payables_invoice ON payables(invoice_number);

-- Assets table (for investments)
CREATE TABLE assets (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    symbol VARCHAR(50) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    type VARCHAR(50) NOT NULL CHECK (type IN ('STOCK', 'BOND', 'MUTUAL_FUND', 'ETF', 'CRYPTO', 'REAL_ESTATE', 'COMMODITY', 'OTHER')),
    currency CHAR(3) NOT NULL DEFAULT 'USD',
    exchange VARCHAR(100),
    isin VARCHAR(12),
    cusip VARCHAR(9),
    sector VARCHAR(100),
    industry VARCHAR(100),
    description TEXT,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    version INTEGER NOT NULL DEFAULT 1
);

CREATE INDEX idx_assets_type ON assets(type);
CREATE INDEX idx_assets_is_active ON assets(is_active);

-- Investment transactions table
CREATE TABLE investment_transactions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    account_id UUID NOT NULL REFERENCES accounts(id) ON DELETE RESTRICT,
    asset_id UUID NOT NULL REFERENCES assets(id) ON DELETE RESTRICT,
    type VARCHAR(50) NOT NULL CHECK (type IN ('BUY', 'SELL', 'DIVIDEND', 'SPLIT', 'MERGE', 'SPIN_OFF')),
    quantity NUMERIC(19, 8) NOT NULL,
    price_per_unit NUMERIC(19, 4) NOT NULL,
    fees NUMERIC(19, 4) NOT NULL DEFAULT 0,
    taxes NUMERIC(19, 4) NOT NULL DEFAULT 0,
    total_amount NUMERIC(19, 4) NOT NULL,
    currency CHAR(3) NOT NULL DEFAULT 'USD',
    transaction_date TIMESTAMPTZ NOT NULL,
    settlement_date TIMESTAMPTZ NOT NULL,
    notes TEXT,
    reference VARCHAR(255),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    version INTEGER NOT NULL DEFAULT 1
);

CREATE INDEX idx_investment_transactions_account ON investment_transactions(account_id);
CREATE INDEX idx_investment_transactions_asset ON investment_transactions(asset_id);
CREATE INDEX idx_investment_transactions_type ON investment_transactions(type);
CREATE INDEX idx_investment_transactions_date ON investment_transactions(transaction_date);

-- Holdings table
CREATE TABLE holdings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    account_id UUID NOT NULL REFERENCES accounts(id) ON DELETE RESTRICT,
    asset_id UUID NOT NULL REFERENCES assets(id) ON DELETE RESTRICT,
    quantity NUMERIC(19, 8) NOT NULL DEFAULT 0,
    average_cost NUMERIC(19, 4) NOT NULL DEFAULT 0,
    total_cost NUMERIC(19, 4) NOT NULL DEFAULT 0,
    current_price NUMERIC(19, 4) NOT NULL DEFAULT 0,
    current_value NUMERIC(19, 4) NOT NULL DEFAULT 0,
    unrealized_gain NUMERIC(19, 4) NOT NULL DEFAULT 0,
    realized_gain NUMERIC(19, 4) NOT NULL DEFAULT 0,
    currency CHAR(3) NOT NULL DEFAULT 'USD',
    first_purchase_date TIMESTAMPTZ,
    last_transaction_date TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    version INTEGER NOT NULL DEFAULT 1,
    UNIQUE(account_id, asset_id)
);

CREATE INDEX idx_holdings_account ON holdings(account_id);
CREATE INDEX idx_holdings_asset ON holdings(asset_id);

-- Capital gains table
CREATE TABLE capital_gains (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    account_id UUID NOT NULL REFERENCES accounts(id) ON DELETE RESTRICT,
    asset_id UUID NOT NULL REFERENCES assets(id) ON DELETE RESTRICT,
    transaction_id UUID NOT NULL REFERENCES investment_transactions(id) ON DELETE RESTRICT,
    type VARCHAR(50) NOT NULL CHECK (type IN ('SHORT_TERM', 'LONG_TERM')),
    quantity NUMERIC(19, 8) NOT NULL,
    cost_basis NUMERIC(19, 4) NOT NULL,
    proceeds NUMERIC(19, 4) NOT NULL,
    gain_loss NUMERIC(19, 4) NOT NULL,
    holding_period INTEGER NOT NULL,
    is_short_term BOOLEAN NOT NULL,
    tax_year INTEGER NOT NULL,
    transaction_date TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_capital_gains_account ON capital_gains(account_id);
CREATE INDEX idx_capital_gains_asset ON capital_gains(asset_id);
CREATE INDEX idx_capital_gains_tax_year ON capital_gains(tax_year);
CREATE INDEX idx_capital_gains_type ON capital_gains(type);

-- Audit log table
CREATE TABLE audit_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    entity_type VARCHAR(100) NOT NULL,
    entity_id UUID NOT NULL,
    action VARCHAR(50) NOT NULL CHECK (action IN ('CREATE', 'UPDATE', 'DELETE')),
    old_value JSONB,
    new_value JSONB,
    user_id UUID,
    ip_address INET,
    user_agent TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_audit_logs_entity ON audit_logs(entity_type, entity_id);
CREATE INDEX idx_audit_logs_user ON audit_logs(user_id);
CREATE INDEX idx_audit_logs_created ON audit_logs(created_at);

-- Trigger function for updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Apply updated_at triggers
CREATE TRIGGER update_accounts_updated_at BEFORE UPDATE ON accounts FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_categories_updated_at BEFORE UPDATE ON categories FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_parties_updated_at BEFORE UPDATE ON parties FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_transactions_updated_at BEFORE UPDATE ON transactions FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_receivables_updated_at BEFORE UPDATE ON receivables FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_payables_updated_at BEFORE UPDATE ON payables FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_assets_updated_at BEFORE UPDATE ON assets FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_investment_transactions_updated_at BEFORE UPDATE ON investment_transactions FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_holdings_updated_at BEFORE UPDATE ON holdings FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();