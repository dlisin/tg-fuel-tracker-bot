package db

import (
	"github.com/dlisin/tg-fuel-tracker-bot/internal/bot/model"
)

func (db *DB) UpsertUser(user model.User) error {
	_, err := db.Exec(`INSERT INTO users(tg_id,car_make,fuel_type,odo_at_register)
		VALUES(?,?,?,?)
		ON CONFLICT(tg_id) DO UPDATE SET car_make=excluded.car_make, fuel_type=excluded.fuel_type, odo_at_register=excluded.odo_at_register`,
		user.TelegramID)
	return err
}

func (db *DB) GetUserByTelegramID(tgID int64) (model.User, error) {
	user := model.User{}

	row := db.QueryRow(`SELECT id,tg_id,car_make,fuel_type,odo_at_register,created_at FROM users WHERE tg_id=?`, tgID)
	if err := row.Scan(&user.ID, &user.TelegramID, &user.CarMake, &user.FuelType, &user.OdoAtRegister, &user.CreatedAt); err != nil {
		return user, err
	}
	return user, nil
}
