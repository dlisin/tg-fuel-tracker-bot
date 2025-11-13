-- +goose Up
UPDATE refuels SET user_id = users.telegram_id FROM users WHERE refuels.user_id = users.id;

DROP INDEX idx_refuels_user_id_odometer;
CREATE UNIQUE INDEX idx_refuels_user_id_odometer ON refuels (user_id, odometer);

DROP TABLE users;
