-- +goose Up
ALTER TABLE cars ADD COLUMN created_by INTEGER NOT NULL DEFAULT 0;

UPDATE cars AS c
    SET created_by = uc.user_id
FROM user_cars AS uc
WHERE c.id = uc.car_id AND uc.is_owner = 1;

DROP INDEX idx__cars__reg_number;

CREATE UNIQUE INDEX idx__cars__created_by__reg_number ON cars (created_by, reg_number);

-- +goose Down
DROP INDEX IF EXISTS idx__cars__created_by__reg_number;

CREATE UNIQUE INDEX idx__cars__reg_number ON cars (reg_number);

ALTER TABLE cars DROP COLUMN created_by;
