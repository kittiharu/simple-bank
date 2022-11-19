package usecase

import (
	"context"

	db "github.com/techschool/simplebank/db/sqlc"
	"github.com/techschool/simplebank/domain"
)

type TransferUsecase interface {
	TransferBalance(context.Context, domain.TransferRequest) (db.TransferTxResult, error)
}

type transferUsecase struct {
	store db.Store
}

func NewTransferUsecase(store db.Store) TransferUsecase {
	return &transferUsecase{
		store: store,
	}
}

func (uc *transferUsecase) TransferBalance(ctx context.Context, req domain.TransferRequest) (db.TransferTxResult, error) {

	arg := db.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}
	result, err := uc.store.TransferTx(ctx, arg)
	return result, err
}
