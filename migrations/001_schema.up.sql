CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    username TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    is_verified BOOLEAN NOT NULL DEFAULT FALSE,
    email_verified_at TIMESTAMP NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP DEFAULT NULL
);

CREATE TABLE email_verification_tokens (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash TEXT NOT NULL UNIQUE,
    expires_at TIMESTAMP NOT NULL,
    used_at TIMESTAMP,
    revoked_at TIMESTAMP NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_evt_user_active
    ON email_verification_tokens (user_id)
    WHERE used_at IS NULL AND revoked_at IS NULL;


CREATE TABLE password_reset_otps (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    otp_hash TEXT NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    used_at TIMESTAMP NULL,
    attempt_count INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_password_reset_otps_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
);

CREATE INDEX idx_password_reset_otps_user_created
    ON password_reset_otps (user_id, created_at DESC);

CREATE INDEX idx_password_reset_otps_expires
    ON password_reset_otps (expires_at);

CREATE TABLE password_reset_tokens (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,

    token_hash TEXT NOT NULL,

    expires_at TIMESTAMP NOT NULL,
    used_at TIMESTAMP NULL,

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_password_reset_tokens_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
);

CREATE UNIQUE INDEX idx_password_reset_tokens_hash
    ON password_reset_tokens (token_hash);

CREATE INDEX idx_password_reset_tokens_expires
    ON password_reset_tokens (expires_at);

CREATE TABLE wallets (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL UNIQUE REFERENCES users(id),
    balance BIGINT NOT NULL DEFAULT 0,
    currency VARCHAR(10) NOT NULL DEFAULT 'IDR',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()  
);

CREATE TYPE wallet_transaction_type AS ENUM (
    'TOPUP',
    'TRANSFER_IN',
    'TRANSFER_OUT',
    'PAYMENT'
);

CREATE TABLE wallet_transactions (
    id BIGSERIAL PRIMARY KEY,
    wallet_id BIGINT NOT NULL REFERENCES wallets(id),
    type wallet_transaction_type NOT NULL,
    status VARCHAR(20) DEFAULT 'SUCCESS',
    idempotency_key VARCHAR(100),
    amount BIGINT NOT NULL,
    balance_before BIGINT NOT NULL,
    balance_after BIGINT NOT NULL,
    reference_id VARCHAR(100),
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE UNIQUE INDEX uniq_wallet_tx_idempotency
    ON wallet_transactions (idempotency_key);

