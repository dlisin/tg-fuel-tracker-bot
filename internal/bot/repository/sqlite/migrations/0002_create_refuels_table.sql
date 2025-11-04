-- +goose Up
CREATE TABLE refuels (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id         INTEGER   NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    odometer        INTEGER   NOT NULL,
    liters          REAL      NOT NULL,
    price_per_liter REAL      NOT NULL,
    price_total     REAL      NOT NULL,
    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_refuels_user_id_created_at ON refuels (user_id, created_at);
CREATE INDEX idx_refuels_user_id_odometer ON refuels (user_id, odometer);

-- +goose Down
DROP INDEX IF EXISTS idx_refuels_user_id_created_at;
DROP INDEX IF EXISTS idx_refuels_user_id_odometer;
DROP TABLE IF EXISTS refuels;
