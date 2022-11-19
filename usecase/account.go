package usecase

import (
	"context"

	db "github.com/techschool/simplebank/db/sqlc"
	"github.com/techschool/simplebank/domain"
)

type AccountUseCase interface {
	CreateAccount(context.Context, domain.CreateAccountRequest, string) (domain.Account, error)
	GetAccount(context.Context, int64) (domain.Account, error)
	ListAccounts(context.Context, domain.ListAccountRequest, string) ([]domain.Account, error)
}

type accountUseCase struct {
	store db.Store
}

func NewAccountUseCase(store db.Store) *accountUseCase {
	return &accountUseCase{
		store: store,
	}
}

func (uc *accountUseCase) CreateAccount(ctx context.Context, req domain.CreateAccountRequest, username string) (domain.Account, error) {
	account, err := uc.store.CreateAccount(ctx, db.CreateAccountParams{
		Owner:    username,
		Currency: req.Currency,
		Amount:   0,
	})

	return domain.NewAccount(account), err
}

func (uc *accountUseCase) GetAccount(ctx context.Context, accountId int64) (domain.Account, error) {
	db, err := uc.store.GetAccount(ctx, accountId)
	account := domain.NewAccount(db)

	return account, err
}

func (uc *accountUseCase) ListAccounts(ctx context.Context, req domain.ListAccountRequest, username string) ([]domain.Account, error) {
	result := make([]domain.Account, 0)
	arg := db.ListAccountsParams{
		Owner:  username,
		Limit:  req.PageSize,
		Offset: (req.PageId - 1) * req.PageSize,
	}
	accounts, err := uc.store.ListAccounts(ctx, arg)
	for _, a := range accounts {
		result = append(result, domain.NewAccount(a))
	}

	return result, err
}
