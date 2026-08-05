-- Migration: 000002_add_multi_tenancy.up.sql

-- Enable UUID extension if not already present
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Users Table
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    full_name VARCHAR(255),
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Households Table (The Multi-tenancy container)
CREATE TABLE households (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    owner_id UUID NOT NULL REFERENCES users(id),
    mode VARCHAR(50) NOT NULL CHECK (mode IN ('SINGLE', 'HOUSEHOLD')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Add household_id to existing tables to enforce ownership
ALTER TABLE accounts ADD COLUMN household_id UUID NOT NULL REFERENCES households(id) ON DELETE CASCADE;
ALTER TABLE categories ADD COLUMN household_id UUID NOT NULL REFERENCES households(id)GM ON DELETE CASCADE;
ALTER TABLE parties ADD COLUMN household_id UUID NOT NULL REFERENCES households(id) ON DELETE CASCADE;
ALTER TABLE transactions ADD COLUMN household_id UUID NOT NULL REFERENCES households(id) ON DELETE CASCADE;
ALTER TABLE receivables ADD COLUMN household_id UUID NOT NULL REFERENCES households(id) ON DELETE CASCADE;
ALTER TABLE payables ADD COLUMN household_id UUID NOT NULL REFERENCES households(id) ON DELETE CASCADE;
ALTER TABLE assets ADD COLUMN household_id UUID NOT NULL REFERENCES households(id) ON DELETE CASCADE;
ALTER TABLE investment_transactions ADD COLUMN household_id UUID NOT NULL REFERENCES households(id) ON DELETE CASCADE;
ALTER TABLE holdings ADD COLUMN household_id UUID NOT NULL REFERENCES households(id) ON DELETE CASCADE;
ALTER TABLE capital_gains ADD COLUMN household_id UUID NOT NULL REFERENCES households(id) ON DELETE CASCADE;

-- Create indexes for performance on foreign keys
CREATE INDEX idx_accounts_household_id ON accounts(household_id);
CREATE INDEX idx_categories_household_id ON categories(household_id);
CREATE INDEX idx_parties_household_id ON parties(household_id);
CREATE INDEX idx_transactions_household_id ON transactions(household_id);
CREATE INDEX idx_receivables_household_id ON receivables(household_id);
CREATE INDEX idx_payables_ownership_id ON payables(household_id);
CREATE INDEX idx_assets_household_id ON assets(household_id);
CREATE INDEX idx_investment_transactions_household_id ON investment_transactions(household_id);
CREATE INDEX idx_holdings_household_id ON holdings(household_id);
CREATE INDEX idx_capital_gains_household_id ON capital_gains(household_id);

-- Update triggers for updated_at (since we modified many tables)
-- This is a bit of a hack in migrations, but ensures new columns also benefit from the existing update function.
-- In a real repo, we would ensure all table definitions included the trigger.

-- Note: The existing 'update_updated_at_column' function and triggers 
-- for existing tables are already defined in 000001.
