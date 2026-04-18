CREATE TABLE IF NOT EXISTS users (
    id           UUID         PRIMARY KEY,
    first_name   VARCHAR(64)  NOT NULL,
    last_name    VARCHAR(64)  NOT NULL,
    phone_number VARCHAR(16)  NOT NULL,
    address      VARCHAR(255) NOT NULL,
    pin_hash     VARCHAR(72)  NOT NULL,
    created_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS users_phone_number_key ON users (phone_number);