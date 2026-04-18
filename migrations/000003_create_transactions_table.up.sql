DO $$ BEGIN
    CREATE TYPE transaction_type AS ENUM ('TOPUP', 'PAYMENT', 'TRANSFER_OUT', 'TRANSFER_IN');
EXCEPTION
    WHEN duplicate_object THEN NULL;
END $$;

DO $$ BEGIN
    CREATE TYPE transaction_status AS ENUM ('PENDING', 'SUCCESS', 'FAILED');
EXCEPTION
    WHEN duplicate_object THEN NULL;
END $$;

CREATE TABLE IF NOT EXISTS transactions (
    id                    UUID               PRIMARY KEY,
    user_id               UUID               NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    type                  transaction_type   NOT NULL,
    status                transaction_status NOT NULL DEFAULT 'PENDING',
    amount                BIGINT             NOT NULL CHECK (amount > 0),
    balance_before        BIGINT             NOT NULL,
    balance_after         BIGINT             NOT NULL,
    remarks               VARCHAR(255)       NOT NULL DEFAULT '',
    reference_id          UUID,
    counterparty_user_id  UUID REFERENCES users (id) ON DELETE SET NULL,
    created_at            TIMESTAMPTZ        NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS transactions_user_id_created_at_idx
    ON transactions (user_id, created_at DESC);

CREATE INDEX IF NOT EXISTS transactions_reference_id_idx
    ON transactions (reference_id);