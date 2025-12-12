-- +goose Up
CREATE TABLE IF NOT EXISTS telegram_integrations (
    id BIGSERIAL PRIMARY KEY,
    shop_id BIGINT NOT NULL UNIQUE,
    bot_token TEXT NOT NULL,
    chat_id TEXT NOT NULL,
    enabled BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_shop FOREIGN KEY (shop_id) REFERENCES shops(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_telegram_integrations_shop_id ON telegram_integrations(shop_id);

-- +goose Down
DROP TABLE IF EXISTS telegram_integrations;


