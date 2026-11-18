-- Enable pgcrypto/uuid if needed (gen_random_uuid is built-in in Postgres 13+)

CREATE TABLE IF NOT EXISTS merchants (
    id BIGINT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    balance BIGINT NOT NULL CHECK (balance >= 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS transfers (
    id UUID PRIMARY KEY,
    merchant_id BIGINT NOT NULL REFERENCES merchants(id),
    idempotency_key VARCHAR(100) NOT NULL,
    partner_id VARCHAR(50) NOT NULL,
    bank_code VARCHAR(20) NOT NULL,
    account_number VARCHAR(50) NOT NULL,
    amount BIGINT NOT NULL CHECK (amount > 0),
    status VARCHAR(20) NOT NULL,
    retry_count INT NOT NULL DEFAULT 0,
    failure_reason TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_transfers_merchant_idempotency
    ON transfers (merchant_id, idempotency_key);

CREATE INDEX IF NOT EXISTS idx_transfers_status_created_at
    ON transfers (status, created_at);

CREATE TABLE IF NOT EXISTS balance_mutations (
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    merchant_id BIGINT NOT NULL REFERENCES merchants(id),
    transfer_id UUID NOT NULL REFERENCES transfers(id),
    type VARCHAR(20) NOT NULL CHECK (type IN ('DEBIT', 'CREDIT')),
    amount BIGINT NOT NULL CHECK (amount > 0),
    balance_before BIGINT NOT NULL,
    balance_after BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_balance_mutations_transfer_id
    ON balance_mutations (transfer_id);

-- Seed 1 hardcoded merchant with initial balance 100,000,000 (100 juta)
INSERT INTO merchants (id, name, balance)
VALUES (1, 'Merchant Utama', 100000000)
ON CONFLICT (id) DO NOTHING;
