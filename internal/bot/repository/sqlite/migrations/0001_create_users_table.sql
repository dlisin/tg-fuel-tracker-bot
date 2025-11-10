-- +goose Up
CREATE TABLE users (
    id               INTEGER PRIMARY KEY AUTOINCREMENT,
    telegram_id      INTEGER   NOT NULL,
    created_at       TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
    
CREATE UNIQUE INDEX idx_telegram_id ON users (telegram_id);

-- +goose Down
DROP INDEX IF EXISTS idx_telegram_id;
DROP TABLE IF EXISTS users;
