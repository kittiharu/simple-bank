package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"github.com/techschool/simplebank/domain"
	"github.com/techschool/simplebank/usecase"
)

type UserHandler struct {
	userUsecase usecase.UserUsecase
}

func NewUserHandler(router *gin.Engine, userUsecase usecase.UserUsecase) {
	handler := &UserHandler{
		userUsecase: userUsecase,
	}
	router.POST("/users", handler.createUser)
	router.POST("/login", handler.login)
}

func (handler *UserHandler) createUser(ctx *gin.Context) {
	var req domain.CreateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := handler.userUsecase.CreateUser(ctx.Request.Context(), req)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "foreign_key_violation", "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (handler *UserHandler) login(ctx *gin.Context) {
	var req domain.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	arg := domain.LoginUserInput{
		LoginRequest: domain.LoginRequest{
			Username: req.Username,
			Password: req.Password,
		},
		UserAgent: ctx.Request.UserAgent(),
		ClientIP:  ctx.ClientIP(),
	}
	res, err := handler.userUsecase.LoginUser(ctx.Request.Context(), arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, res)
}
