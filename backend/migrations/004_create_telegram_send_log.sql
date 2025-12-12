-- +goose Up
CREATE TYPE telegram_send_log_status AS ENUM ('SENT', 'FAILED');

CREATE TABLE IF NOT EXISTS telegram_send_log (
    id BIGSERIAL PRIMARY KEY,
    shop_id BIGINT NOT NULL,
    order_id BIGINT NOT NULL,
    message TEXT NOT NULL,
    status telegram_send_log_status NOT NULL,
    error TEXT,
    sent_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_shop FOREIGN KEY (shop_id) REFERENCES shops(id) ON DELETE CASCADE,
    CONSTRAINT fk_order FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE,
    CONSTRAINT unique_shop_order UNIQUE (shop_id, order_id)
);

CREATE INDEX IF NOT EXISTS idx_telegram_send_log_shop_id ON telegram_send_log(shop_id);
CREATE INDEX IF NOT EXISTS idx_telegram_send_log_order_id ON telegram_send_log(order_id);
CREATE INDEX IF NOT EXISTS idx_telegram_send_log_sent_at ON telegram_send_log(sent_at);

-- +goose Down
DROP TABLE IF EXISTS telegram_send_log;
DROP TYPE IF EXISTS telegram_send_log_status;


