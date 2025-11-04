-- +goose Up
CREATE TABLE users (
    id               INTEGER PRIMARY KEY AUTOINCREMENT,
    telegram_id      INTEGER   NOT NULL,
    fuel_type        TEXT      NOT NULL,
    currency         TEXT      NOT NULL,
    current_odometer INTEGER   NOT NULL DEFAULT 0,
    created_at       TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at       TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX idx_telegram_id ON users (telegram_id);

-- +goose Down
DROP INDEX IF EXISTS idx_telegram_id;
DROP TABLE IF EXISTS users;
