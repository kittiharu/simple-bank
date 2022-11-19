package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/techschool/simplebank/db/sqlc"
	"github.com/techschool/simplebank/token"
	"github.com/techschool/simplebank/usecase"
	"github.com/techschool/simplebank/util"
)

type Server struct {
	config     util.Config
	tokenMaker token.Maker
	router     *gin.Engine
}

func NewServer(store db.Store, config util.Config) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:     config,
		tokenMaker: tokenMaker,
	}
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	router := gin.Default()
	authRoutes := router.Group("/").Use(authMiddleware(tokenMaker))

	userUsecase := usecase.NewUserUsecase(store, tokenMaker, config)
	sessionUsecase := usecase.NewSessionUsecase(store)
	accountUseCase := usecase.NewAccountUseCase(store)
	transferUsecase := usecase.NewTransferUsecase(store)

	NewAccountHandler(server, authRoutes, accountUseCase)
	NewTransferHandler(authRoutes, transferUsecase, accountUseCase)
	NewUserHandler(router, userUsecase)
	NewTokenHandler(router, sessionUsecase, tokenMaker, config)

	server.router = router
	return server, nil
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func (server *Server) Run() {
	server.router.Run(server.config.HttpServerAddress)
}
