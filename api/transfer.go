package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/techschool/simplebank/domain"
	"github.com/techschool/simplebank/token"
	"github.com/techschool/simplebank/usecase"
)

type TransferHandler struct {
	transferUsecase usecase.TransferUsecase
	accountUsecase  usecase.AccountUseCase
}

func NewTransferHandler(routes gin.IRoutes, transferUsecase usecase.TransferUsecase, accountUsecase usecase.AccountUseCase) {
	handler := &TransferHandler{
		transferUsecase: transferUsecase,
		accountUsecase:  accountUsecase,
	}
	routes.POST("/transfer", handler.transferBalance)
}

func (handler *TransferHandler) transferBalance(ctx *gin.Context) {
	var req domain.TransferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	fromAccount, valid := handler.validAccount(ctx, req.FromAccountID, req.Currency)
	if !valid {
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if authPayload.Username != fromAccount.Owner {
		err := errors.New("from accont doesn't belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	_, valid = handler.validAccount(ctx, req.ToAccountID, req.Currency)
	if !valid {
		return
	}

	result, err := handler.transferUsecase.TransferBalance(ctx.Request.Context(), req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (handler *TransferHandler) validAccount(ctx *gin.Context, accountId int64, currency string) (domain.Account, bool) {
	account, err := handler.accountUsecase.GetAccount(ctx, accountId)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return account, false
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return account, false
	}

	if account.Currency != currency {
		err := fmt.Errorf("account id %d currency mismatch %s vs %s", accountId, account.Currency, currency)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return account, false
	}

	return account, true
}
