CREATE TABLE IF NOT EXISTS products (
    id          BIGINT PRIMARY KEY,
    sku         VARCHAR(64) NOT NULL UNIQUE,
    name        VARCHAR(255) NOT NULL,
    qty         INT NOT NULL CHECK (qty >= 0),
    price       NUMERIC(18,2) NOT NULL CHECK (price > 0),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_products_sku ON products (sku);
