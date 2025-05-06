package user

import "time"

type UserRepo struct {
	ID           int    `json:"id"`
	TgID         int64  `json:"tg_id"`
	Username     string `json:"username"`
	Balance      int64  `json:"balance"`
	Level        int    `json:"level"`
	Energy       int    `json:"energy"`
	MaxEnergy    int    `json:"max_energy"`
	ProfitPerHour int    `json:"profit_per_hour"`
	Head         *string `json:"head"`
	Body         *string `json:"body"`
	Legs         *string `json:"legs"`
	Foot         *string `json:"foot"`
	Hand         *string `json:"hand"`
	ProfitForTap int    `json:"profit_for_tap"`
	LastRestoration time.Time `json:"last_restoration"`
	LastProfitPerHour time.Time `json:"last_profit_per_hour"`
}

type UserResponse struct {
	Username     string `json:"username"`
	Balance      int64  `json:"balance"`
	Level        int    `json:"level"`
	Energy       int    `json:"energy"`
	MaxEnergy    int    `json:"max_energy"`
	ProfitPerHour int    `json:"profit_per_hour"`
	Head         *string `json:"head"`
	Body         *string `json:"body"`
	Legs         *string `json:"legs"`
	Foot         *string `json:"foot"`
	Hand         *string `json:"hand"`
	ProfitForTap int    `json:"profit_for_tap"`
}

type UserRequest struct {
	Username string `json:"username"`
	Balance  int    `json:"balance"`
	Level    int    `json:"level"`
	CoinForTap         int    `json:"coin_for_tap"`
	Energy             int    `json:"energy"`
	ProfitPerHour      int    `json:"profit_per_hour"`
}

type CreateUserRequest struct {
	TgID     int64  `json:"tg_id"`
	Username string `json:"username"`
}

type TgIDRequest struct {
	TgID int64 `header:"tg_id"`
}



