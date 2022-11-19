package api

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/techschool/simplebank/domain"
	"github.com/techschool/simplebank/token"
	"github.com/techschool/simplebank/usecase"
)

type AccountHandler struct {
	*Server
	accountUseCase usecase.AccountUseCase
}

func NewAccountHandler(server *Server, routes gin.IRoutes, accountUseCase usecase.AccountUseCase) {
	handler := &AccountHandler{
		accountUseCase: accountUseCase,
		Server:         server,
	}
	routes.POST("/accounts", handler.createAccount)
	routes.GET("/accounts/:id", handler.getAccount)
	routes.GET("/accounts", handler.listAccounts)
}

func (handler *AccountHandler) createAccount(ctx *gin.Context) {
	var req domain.CreateAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	account, err := handler.accountUseCase.CreateAccount(ctx.Request.Context(), req, authPayload.Username)

	if err != nil {
		// if pqErr, ok := err.(*pq.Error); ok {
		// 	switch pqErr.Code.Name() {
		// 	case "foreign_key_violation", "unique_violation":
		// 		ctx.JSON(http.StatusForbidden, errorResponse(err))
		// 		return
		// 	}
		// }
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

func (handler *AccountHandler) getAccount(ctx *gin.Context) {
	var req domain.GetAccountRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	account, err := handler.accountUseCase.GetAccount(ctx.Request.Context(), req.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if account.Owner != authPayload.Username {
		err := errors.New("account doesn't belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

func (handler *AccountHandler) listAccounts(ctx *gin.Context) {
	var req domain.ListAccountRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	account, err := handler.accountUseCase.ListAccounts(ctx.Request.Context(), req, authPayload.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}
