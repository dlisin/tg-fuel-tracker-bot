-- +goose Up
CREATE TABLE cars (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    reg_number TEXT NOT NULL,
    fuel_type TEXT NOT NULL,
    odometer INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX idx__cars__reg_number ON cars (reg_number);

INSERT INTO cars (reg_number, fuel_type, odometer, created_at, updated_at)
SELECT
    r.user_id as reg_number,
    "N/A" as fuel_type,
    MAX(r.odometer) as odometer,
    MIN(r.created_at),
    MAX(r.created_at)
FROM refuels r
GROUP BY
    r.user_id;

-- +goose Down
DROP INDEX IF EXISTS idx__cars__reg_number;
DROP TABLE IF EXISTS cars;
