-- +goose Up
CREATE TABLE user_car_invites (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    car_id INTEGER NOT NULL REFERENCES cars (id) ON DELETE CASCADE,
    token TEXT NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_by INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX idx__user_car_invites__token ON user_car_invites (token);
CREATE INDEX idx__user_car_invites__car_id ON user_car_invites(car_id);

-- +goose Down
DROP INDEX IF EXISTS idx__user_car_invites__car_id;
DROP INDEX IF EXISTS idx__user_car_invites__token;
DROP TABLE IF EXISTS user_car_invites;
