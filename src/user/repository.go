package user

import (
	"database/sql"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetUser(tg_id int64) (*UserRepo, error) {
	var user UserRepo
	var head, body, legs, foot sql.NullString
	err := r.db.QueryRow(
		"SELECT id, tg_id, username, balance, level, energy, max_energy, profit_per_hour, head, body, legs, foot, profit_for_tap, last_restoration, last_profit_per_hour FROM users WHERE tg_id = $1", 
		tg_id,
	).Scan(
		&user.ID, &user.TgID, &user.Username, &user.Balance, &user.Level, &user.Energy, &user.MaxEnergy, &user.ProfitPerHour, &head, &body, &legs, &foot, &user.ProfitForTap, &user.LastRestoration, &user.LastProfitPerHour,
	)
	if err != nil {
		return nil, err
	}
	
	if head.Valid {
		user.Head = &head.String
	}
	if body.Valid {
		user.Body = &body.String
	}
	if legs.Valid {
		user.Legs = &legs.String
	}
	if foot.Valid {
		user.Foot = &foot.String
	}
	
	return &user, nil
}

func (r *UserRepository) CreateUser(tg_id int64, username string) (*UserRepo, error) {
	var user UserRepo
	err := r.db.QueryRow(
		`INSERT INTO users (tg_id, username) 
		VALUES ($1, $2) 
		RETURNING id, tg_id, username, balance, level, energy, max_energy, profit_per_hour, head, body, legs, last_profit_per_hour`,
		tg_id, username,
	).Scan(
		&user.ID, &user.TgID, &user.Username, &user.Balance, &user.Level, &user.Energy, &user.MaxEnergy, &user.ProfitPerHour, &user.Head, &user.Body, &user.Legs, &user.LastProfitPerHour,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) SelectProfitForTap(tg_id int64) (int, error) {
	var profit int
	err := r.db.QueryRow("SELECT profit_for_tap FROM users WHERE tg_id = $1", tg_id).Scan(&profit)
	if err != nil {
		return 0, err
	}
	if profit <= 0 {
		return 1, nil
	}
	return profit, nil
}

func (r *UserRepository) UpdateBalanceForTap(tg_id int64, balance int) error {
	_, err := r.db.Exec("UPDATE users SET balance = balance + $1, energy = energy - 1 WHERE tg_id = $2", balance, tg_id)
	return err
}

func (r *UserRepository) UpdateEnergy(tg_id int64, energy int) error {
	_, err := r.db.Exec("UPDATE users SET energy = energy + $1, last_restoration = NOW() WHERE tg_id = $2", energy, tg_id)
	return err
}

func (r *UserRepository) UpdateBalanceForProfitPerHour(tg_id int64, balance int) error {
	_, err := r.db.Exec("UPDATE users SET balance = $1, last_profit_per_hour = NOW() WHERE tg_id = $2", balance, tg_id)
	return err
}