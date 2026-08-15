-- +goose Up
PRAGMA foreign_keys = 0;
ALTER TABLE refuels ADD COLUMN car_id INTEGER NOT NULL DEFAULT 0 REFERENCES cars(id) ON DELETE CASCADE;

CREATE UNIQUE INDEX idx__refuels__car_id__odometer ON refuels (car_id, odometer);

DROP INDEX idx_refuels_user_id_created_at;
DROP INDEX idx_refuels_user_id_odometer;

UPDATE refuels AS r
SET
    car_id = uc.car_id
FROM user_cars AS uc
WHERE r.user_id = uc.user_id;

PRAGMA foreign_keys = 1;

ALTER TABLE refuels RENAME COLUMN user_id TO created_by;

-- +goose Down
DROP INDEX IF EXISTS idx__refuels__car_id__odometer;

ALTER TABLE refuels RENAME COLUMN created_by TO user_id;
ALTER TABLE refuels DROP COLUMN car_id;

CREATE INDEX idx_refuels_user_id_created_at ON refuels (user_id, created_at);
CREATE INDEX idx_refuels_user_id_odometer ON refuels (user_id, odometer);
