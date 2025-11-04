package db

import (
	"database/sql"
	"embed"
	"errors"
	"time"

	"github.com/pressly/goose/v3"
	_ "modernc.org/sqlite"
)

type DB struct {
	*sql.DB
}

func Open(path string) (sql.DB, error) {
	d, err := sql.Open("sqlite", path+"?_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)")
	if err != nil {
		return nil, err
	}
	if err := migrate(d); err != nil {
		return nil, err
	}
	return &DB{d}, nil
}

func migrate(db *sql.DB) error {
	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("sqlite"); err != nil {
		return err
	}

	if err := goose.Up(db, "migrations"); err != nil {
		return err
	}

	return nil
}

func (db *DB) GetLastOdometer(userID int64) (int64, error) {
	var odo sql.NullInt64
	if err := db.QueryRow(`SELECT odometer FROM fillups WHERE user_id=? ORDER BY odometer DESC LIMIT 1`, userID).Scan(&odo); err != nil && err != sql.ErrNoRows {
		return 0, err
	}
	if odo.Valid {
		return odo.Int64, nil
	}
	return 0, nil
}

func (db *DB) AddFillup(userID int64, odometer int64, liters float64, totalPrice *float64, pricePerLiter *float64) error {
	if liters <= 0 {
		return errors.New("литры должны быть > 0")
	}
	var ppl float64
	if pricePerLiter != nil {
		if *pricePerLiter <= 0 {
			return errors.New("цена/л должна быть > 0")
		}
		ppl = *pricePerLiter
	} else if totalPrice != nil {
		if *totalPrice <= 0 {
			return errors.New("сумма чека должна быть > 0")
		}
		ppl = *totalPrice / liters
	} else {
		return errors.New("укажите сумму чека или цену за литр")
	}
	var tp sql.NullFloat64
	if totalPrice != nil {
		tp = sql.NullFloat64{Float64: *totalPrice, Valid: true}
	}
	_, err := db.Exec(`INSERT INTO fillups(user_id, odometer, liters, total_price, price_per_liter) VALUES(?,?,?,?,?)`,
		userID, odometer, liters, tp, ppl)
	return err
}

func (db *DB) GetFillups(userID int64, startEnd *[2]time.Time) ([]Fillup, error) {
	q := `SELECT id, user_id, ts, odometer, liters, total_price, price_per_liter FROM fillups WHERE user_id=?`
	var rows *sql.Rows
	var err error
	if startEnd != nil {
		q += ` AND ts >= ? AND ts < ? ORDER BY ts ASC`
		rows, err = db.Query(q, userID, startEnd[0], startEnd[1])
	} else {
		q += ` ORDER BY ts ASC`
		rows, err = db.Query(q, userID)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Fillup
	for rows.Next() {
		var f Fillup
		if err := rows.Scan(&f.ID, &f.UserID, &f.TS, &f.Odometer, &f.Liters, &f.TotalPrice, &f.PricePerLiter); err != nil {
			return nil, err
		}
		out = append(out, f)
	}
	return out, rows.Err()
}
