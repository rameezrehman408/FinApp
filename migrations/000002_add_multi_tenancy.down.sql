-- Migration: 000002_add_multi_tenancy.down.sql

-- Remove foreign keys and indexes first (to avoid dependency errors)
DROP INDEX IF EXISTS idx_capital_gains_household_id;
DROP INDEX IF EXISTS idx_holdings_household_id;
DROP INDEX IF EXISTS idx_investment_transactions_household_id;
DROP INDEX IF EXISTS idx_assets_household_id;
DROP INDEX IF EXISTS idx_payables_ownership_id;
DROP INDEX IF EXISTS idx_receivables_household_id;
DROP INDEX IF EXISTS idx_transactions_household_id;
DROP INDEX IF EXISTS idx_parties_household_id;
DROP INDEX IF EXISTS idx_categories_household_id;
DROP INDEX IF EXISTS idx_accounts_household_id;

-- Remove columns from tables
ALTER TABLE capital_gains DROP COLUMN household_id;
ALTER TABLE holdings DROP COLUMN household_id;
ALTER TABLE investment_transactions DROP COLUMN household_id;
ALTER TABLE assets DROP COLUMN household_id;
ALTER TABLE payables DROP COLUMN household_id;
ALTER TABLE receivables DROP COLUMN household_id;
ALTER TABLE transactions DROP COLUMN household_id;
ALTER TABLE parties DROP COLUMN household_id;
ALTER TABLE categories DROP COLUMN household_id;
ALTER TABLE accounts DROP COLUMN household_id;

-- Drop tables (order matters due to FK dependencies)
DROP TABLE IF EXISTS households;
DROP TABLE IF EXISTS users;
