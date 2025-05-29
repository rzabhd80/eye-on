CREATE TABLE order_histories
(
    id                     UUID PRIMARY KEY        DEFAULT gen_random_uuid(),
    user_id                UUID           NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    exchange_credential_id UUID           NOT NULL REFERENCES exchange_credentials (id) ON DELETE CASCADE,
    exchange_id            UUID           NOT NULL REFERENCES exchanges (id) ON DELETE CASCADE,
    trading_pair_id        UUID           NOT NULL REFERENCES trading_pairs (id) ON DELETE RESTRICT,
    client_order_id        VARCHAR(100)   NOT NULL, -- your idempotency key
    exchange_order_id      VARCHAR(100)   NOT NULL, -- exchange’s own order ID
    side                   VARCHAR(10)    NOT NULL, -- buy / sell
    type                   VARCHAR(10)    NOT NULL, -- limit / market
    quantity               DECIMAL(20, 8) NOT NULL,
    price                  DECIMAL(20, 8),
    status                 VARCHAR(20)    NOT NULL,
    created_at             TIMESTAMPTZ    NOT NULL DEFAULT now(),
    updated_at             TIMESTAMPTZ    NOT NULL DEFAULT now(),
    deleted_at             TIMESTAMPTZ
);
-- Composite index for “open orders” lookup
CREATE INDEX idx_orders_user_ex_cred_status
    ON order_histories (user_id, exchange_credential_id, status);
CREATE INDEX idx_order_histories_order_id
    ON order_histories (exchange_order_id);
CREATE INDEX idx_order_histories_client_order_id
    ON order_histories (client_order_id);
