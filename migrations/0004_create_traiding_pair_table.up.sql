CREATE TABLE trading_pairs
(
    id           UUID PRIMARY KEY     DEFAULT gen_random_uuid(),
    exchange_id  UUID        NOT NULL REFERENCES exchanges (id) ON DELETE CASCADE,
    UNIQUE (exchange_id, id),
    symbol       VARCHAR(20) NOT NULL,
    base_asset   VARCHAR(10) NOT NULL,
    quote_asset  VARCHAR(10) NOT NULL,
    min_quantity DECIMAL(20, 8),
    max_quantity DECIMAL(20, 8),
    step_size    DECIMAL(20, 8),
    tick_size    DECIMAL(20, 8),
    is_active    BOOLEAN     NOT NULL DEFAULT true,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at   TIMESTAMPTZ
);
CREATE UNIQUE INDEX ux_trading_pairs_exchange_symbol_active
    ON trading_pairs (exchange_id, symbol) WHERE deleted_at IS NULL;
CREATE INDEX idx_trading_pairs_symbol ON trading_pairs (symbol);
