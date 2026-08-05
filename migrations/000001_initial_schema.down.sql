-- Down migration for initial schema
-- Migration: 000001_initial_schema.down.sql

DROP TRIGGER IF EXISTS update_holdings_updated_at ON holdings;
DROP TRIGGER IF EXISTS update_investment_transactions_updated_at ON investment_transactions;
DROP TRIGGER IF EXISTS update_assets_updated_at ON assets;
DROP TRIGGER IF EXISTS update_payables_updated_at ON payables;
DROP TRIGGER IF EXISTS update_receivables_updated_at ON receivables;
DROP TRIGGER IF EXISTS update_transactions_updated_at ON transactions;
DROP TRIGGER IF EXISTS update_parties_updated_at ON parties;
DROP TRIGGER IF EXISTS update_categories_updated_at ON categories;
DROP TRIGGER IF EXISTS update_accounts_updated_at ON accounts;

DROP FUNCTION IF EXISTS update_updated_at_column();

DROP TABLE IF EXISTS audit_logs;
DROP TABLE IF EXISTS capital_gains;
DROP TABLE IF EXISTS holdings;
DROP TABLE IF EXISTS investment_transactions;
DROP TABLE IF EXISTS assets;
DROP TABLE IF EXISTS payables;
DROP TABLE IF EXISTS receivables;
DROP TABLE IF EXISTS transactions;
DROP TABLE IF EXISTS parties;
DROP TABLE IF EXISTS categories;
DROP TABLE IF EXISTS accounts;

DROP EXTENSION IF EXISTS "uuid-ossp";