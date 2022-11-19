package usecase

import (
	"context"

	"github.com/google/uuid"
	db "github.com/techschool/simplebank/db/sqlc"
)

type SessionUsecase interface {
	GetSession(context.Context, uuid.UUID) (db.Session, error)
}

type sessionUsecase struct {
	store db.Store
}

func NewSessionUsecase(store db.Store) SessionUsecase {
	return &sessionUsecase{
		store: store,
	}
}

func (uc *sessionUsecase) GetSession(ctx context.Context, id uuid.UUID) (db.Session, error) {
	session, err := uc.store.GetSession(ctx, id)
	return session, err
}
