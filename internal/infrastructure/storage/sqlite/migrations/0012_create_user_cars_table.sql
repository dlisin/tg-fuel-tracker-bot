-- +goose Up
CREATE TABLE user_cars (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    car_id INTEGER NOT NULL REFERENCES cars (id) ON DELETE CASCADE,
    is_owner BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX idx__user_cars__user_id__car_id ON user_cars (user_id, car_id);
CREATE UNIQUE INDEX idx__user_cars__car_id__is_owner ON user_cars(car_id) WHERE is_owner = 1;

INSERT INTO user_cars (user_id, car_id, is_owner, created_at)
SELECT
    c.reg_number AS user_id,
    c.id AS car_id,
    TRUE as is_owner,
    c.created_at as created_at
FROM
    cars c;

-- +goose Down
DROP INDEX IF EXISTS idx__user_cars__car_id__is_owner;
DROP INDEX IF EXISTS idx__user_cars__user_id__car_id;
DROP TABLE IF EXISTS user_cars;
