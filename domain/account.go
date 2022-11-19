package domain

import (
	"time"

	db "github.com/techschool/simplebank/db/sqlc"
)

type Account struct {
	ID        int64     `json:"id"`
	Owner     string    `json:"owner"`
	Amount    int64     `json:"amount"`
	Currency  string    `json:"currency"`
	CreatedAt time.Time `json:"created_at"`
}

func NewAccount(db db.Account) Account {
	return Account{
		ID:        db.ID,
		Owner:     db.Owner,
		Amount:    db.Amount,
		Currency:  db.Currency,
		CreatedAt: db.CreatedAt,
	}
}

type CreateAccountRequest struct {
	Currency string `json:"currency" binding:"required,oneof=USD THB"`
}

type GetAccountRequest struct {
	Id int64 `uri:"id" binding:"required,min=1"`
}

type ListAccountRequest struct {
	SearchRequest
}
