CREATE TABLE order_book_snapshots
(
    id              BIGSERIAL PRIMARY KEY,
    trading_pair_id UUID        NOT NULL REFERENCES trading_pairs (id) ON DELETE CASCADE,
    bids            JSONB       NOT NULL, -- e.g. [ [price, qty], … ]
    asks            JSONB       NOT NULL,
    snapshot_time   TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_ob_snapshots_pair_time
    ON order_book_snapshots (trading_pair_id, snapshot_time);
