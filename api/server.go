package api

import (
	"fmt"
	db "simplebank/db/sqlc"
	"simplebank/token"
	"simplebank/util"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	store      db.Store
	router     *gin.Engine
	tokenMaker token.TokenMaker
	config     util.Config
}

func NewServer(config util.Config, st db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %v", err)
	}

	server := Server{
		store:      st,
		tokenMaker: tokenMaker,
		config:     config,
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	server.setupRouter()
	return &server, nil
}

func (s *Server) setupRouter() {
	router := gin.Default()
	router.POST("/users", s.createUser)
	router.POST("/users/login", s.loginUser)

	router.POST("/accounts", s.createAccount)
	router.GET("/accounts/:id", s.getAccount)
	router.GET("/accounts", s.listAccounts)

	router.POST("/transfers", s.createTransfer)

	s.router = router
}

func (s *Server) Start(add string) error {
	return s.router.Run(add)
}

func (s *Server) errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
