CREATE TABLE order_book_snapshots
(
    id              BIGSERIAL PRIMARY KEY,
    exchange_id     UUID        NOT NULL REFERENCES exchanges (id) ON DELETE CASCADE,
    trading_pair_id UUID        NOT NULL REFERENCES trading_pairs (id) ON DELETE CASCADE,
    bids            JSONB       NOT NULL,
    asks            JSONB       NOT NULL,
    snapshot_time   TIMESTAMPTZ NOT NULL DEFAULT now(),

    -- Ensure we don't store snapshots for invalid exchange-pair combinations
    CONSTRAINT fk_valid_exchange_pair
        FOREIGN KEY (exchange_id, trading_pair_id)
            REFERENCES trading_pairs (exchange_id, id)
);

CREATE INDEX idx_ob_snapshots_exchange_pair_time
    ON order_book_snapshots (exchange_id, trading_pair_id, snapshot_time);