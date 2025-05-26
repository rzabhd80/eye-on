CREATE TABLE exchanges
(
    id           UUID PRIMARY KEY      DEFAULT gen_random_uuid(),
    name         VARCHAR(50)  NOT NULL,
    display_name VARCHAR(100) NOT NULL,
    base_url     VARCHAR(255) NOT NULL,
    is_active    BOOLEAN      NOT NULL DEFAULT true,
    rate_limit   INTEGER      NOT NULL DEFAULT 1000,
    features     JSONB,
    created_at   TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at   TIMESTAMPTZ  NOT NULL DEFAULT now(),
    deleted_at   TIMESTAMPTZ
);

CREATE UNIQUE INDEX ux_exchanges_name_active ON exchanges (name) WHERE deleted_at IS NULL;
