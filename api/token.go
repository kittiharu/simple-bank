package api

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/techschool/simplebank/domain"
	"github.com/techschool/simplebank/token"
	"github.com/techschool/simplebank/usecase"
	"github.com/techschool/simplebank/util"
)

type TokenHandler struct {
	sessionUsecase usecase.SessionUsecase
	tokenMaker     token.Maker
	config         util.Config
}

func NewTokenHandler(router *gin.Engine, sessionUsecase usecase.SessionUsecase, tokenMaker token.Maker, config util.Config) {
	handler := &TokenHandler{
		sessionUsecase: sessionUsecase,
		tokenMaker:     tokenMaker,
		config:         config,
	}
	router.POST("/tokens/renew_access", handler.renewAccessToken)
}

func (handler *TokenHandler) renewAccessToken(ctx *gin.Context) {
	var req domain.RefreshTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	payload, err := handler.tokenMaker.VerifyToken(req.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	session, err := handler.sessionUsecase.GetSession(ctx, payload.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if session.IsBlocked {
		err := errors.New("session blocked")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	if payload.Username != session.Username {
		err := errors.New("incorrect user session")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	if session.RefreshToken != req.RefreshToken {
		err := errors.New("mismatch token")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	if time.Now().After(session.ExpiresAt) {
		err := errors.New("session expired")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	accessToken, accessPayload, err := handler.tokenMaker.CreateToken(payload.Username, handler.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := domain.RefreshTokenResponse{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessPayload.ExpiredAt,
	}
	ctx.JSON(http.StatusOK, res)
}
