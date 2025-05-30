CREATE TABLE order_events
(
    id           UUID PRIMARY KEY      DEFAULT gen_random_uuid(),
    order_hist_id UUID           NOT NULL REFERENCES order_histories (id) ON DELETE CASCADE,
    event_type    VARCHAR(30)    NOT NULL, -- e.g. 'new', 'partial_fill', 'filled', 'canceled'
    filled_qty    DECIMAL(20, 8) NOT NULL,
    remaining_qty DECIMAL(20, 8) NOT NULL,
    event_time    TIMESTAMPTZ    NOT NULL,
    recorded_at   TIMESTAMPTZ    NOT NULL DEFAULT now(),
    created_at TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);
CREATE INDEX idx_order_events_order_hist_id
    ON order_events (order_hist_id);
