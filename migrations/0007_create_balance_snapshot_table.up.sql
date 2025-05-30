CREATE TABLE balance_snapshots
(
    id            UUID PRIMARY KEY      DEFAULT gen_random_uuid(),
    user_id       UUID            NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    exchange_id   UUID            NOT NULL REFERENCES exchanges (id) ON DELETE CASCADE,
    currency      VARCHAR(10)     NOT NULL,
    total         NUMERIC(30, 10) NOT NULL,
    available     NUMERIC(30, 10) NOT NULL,
    snapshot_time TIMESTAMPTZ     NOT NULL DEFAULT now(),
    created_at TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);
CREATE INDEX idx_balance_snapshots_user_time
    ON balance_snapshots (user_id, snapshot_time);
