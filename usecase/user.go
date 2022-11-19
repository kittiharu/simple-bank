package usecase

import (
	"context"

	db "github.com/techschool/simplebank/db/sqlc"
	"github.com/techschool/simplebank/domain"
	"github.com/techschool/simplebank/token"
	"github.com/techschool/simplebank/util"
)

type UserUsecase interface {
	CreateUser(context.Context, domain.CreateUserRequest) (domain.UserResponse, error)
	LoginUser(context.Context, domain.LoginUserInput) (domain.LoginResponse, error)
}

type userUsecase struct {
	store      db.Store
	tokenMaker token.Maker
	config     util.Config
}

func NewUserUsecase(store db.Store, tokenMaker token.Maker, config util.Config) UserUsecase {
	return &userUsecase{
		store:      store,
		tokenMaker: tokenMaker,
		config:     config,
	}
}

func (uc *userUsecase) CreateUser(ctx context.Context, req domain.CreateUserRequest) (domain.UserResponse, error) {
	hashPassword, err := util.HashPassword(req.Password)
	if err != nil {
		return domain.UserResponse{}, err
	}

	arg := db.CreateUserParams{
		Username:       req.Username,
		FullName:       req.FullName,
		Email:          req.Email,
		HashedPassword: hashPassword,
	}

	user, err := uc.store.CreateUser(ctx, arg)
	res := domain.NewUserResponse(user)

	return res, err
}

func (uc *userUsecase) LoginUser(ctx context.Context, req domain.LoginUserInput) (domain.LoginResponse, error) {
	user, err := uc.store.GetUser(ctx, req.Username)
	if err != nil {
		return domain.LoginResponse{}, err
	}

	err = util.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		return domain.LoginResponse{}, err
	}

	token, payload, err := uc.tokenMaker.CreateToken(user.Username, uc.config.AccessTokenDuration)
	if err != nil {
		return domain.LoginResponse{}, err
	}

	refreshToken, refreshPayload, err := uc.tokenMaker.CreateToken(user.Username, uc.config.RefreshTokenDuration)
	if err != nil {
		return domain.LoginResponse{}, err
	}

	session, err := uc.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshPayload.ID,
		Username:     req.Username,
		RefreshToken: refreshToken,
		UserAgent:    req.UserAgent,
		ClientIp:     req.ClientIP,
		IsBlocked:    false,
		ExpiresAt:    refreshPayload.ExpiredAt,
	})
	if err != nil {
		return domain.LoginResponse{}, err
	}

	res := domain.LoginResponse{
		SessionID:             session.ID,
		AccessToken:           token,
		AccessTokenExpiresAt:  payload.ExpiredAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshPayload.ExpiredAt,
		User:                  domain.NewUserResponse(user),
	}

	return res, err
}
