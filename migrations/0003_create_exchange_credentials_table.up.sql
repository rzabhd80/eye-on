CREATE TABLE exchange_credentials
(
    id          UUID PRIMARY KEY      DEFAULT gen_random_uuid(),
    user_id     UUID         NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    exchange_id UUID         NOT NULL REFERENCES exchanges (id) ON DELETE CASCADE,
    label       VARCHAR(100) NOT NULL DEFAULT 'Default',
    api_key     TEXT         NOT NULL,
    secret_key  TEXT,
    passphrase  TEXT,
    refresh_key TEXT,
    is_active   BOOLEAN      NOT NULL DEFAULT true,
    is_testnet  BOOLEAN      NOT NULL DEFAULT false,
    permissions JSONB,
    last_used   TIMESTAMPTZ,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT now(),
    deleted_at  TIMESTAMPTZ,
);

CREATE INDEX idx_exchange_credentials_user_exchange ON
    exchange_credentials (user_id, exchange_id);
