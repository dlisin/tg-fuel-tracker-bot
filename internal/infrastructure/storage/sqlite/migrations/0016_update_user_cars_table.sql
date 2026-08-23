-- +goose Up
DROP INDEX idx__user_cars__car_id__is_owner;

ALTER TABLE user_cars DROP COLUMN is_owner;

-- +goose Down
ALTER TABLE user_cars ADD COLUMN is_owner BOOLEAN NOT NULL DEFAULT FALSE;

UPDATE user_cars AS uc
    SET is_owner = TRUE
FROM cars AS c
WHERE uc.car_id = c.id AND uc.user_id = c.created_by;

CREATE UNIQUE INDEX idx__user_cars__car_id__is_owner ON user_cars (car_id) WHERE is_owner = 1;
