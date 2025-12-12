-- +goose Up
CREATE TABLE IF NOT EXISTS orders (
    id BIGSERIAL PRIMARY KEY,
    shop_id BIGINT NOT NULL,
    number TEXT NOT NULL,
    total NUMERIC(10, 2) NOT NULL,
    customer_name TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_shop FOREIGN KEY (shop_id) REFERENCES shops(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_orders_shop_id ON orders(shop_id);

-- +goose Down
DROP TABLE IF EXISTS orders;


