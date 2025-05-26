CREATE TABLE balance_snapshots
(
    id            BIGSERIAL PRIMARY KEY,
    user_id       UUID            NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    exchange_id   UUID            NOT NULL REFERENCES exchanges (id) ON DELETE CASCADE,
    currency      VARCHAR(10)     NOT NULL,
    total         NUMERIC(30, 10) NOT NULL,
    available     NUMERIC(30, 10) NOT NULL,
    snapshot_time TIMESTAMPTZ     NOT NULL DEFAULT now()
);
CREATE INDEX idx_balance_snapshots_user_time
    ON balance_snapshots (user_id, snapshot_time);
