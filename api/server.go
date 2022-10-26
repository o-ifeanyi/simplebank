package api

import (
	db "simplebank/db/sqlc"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	store  db.Store
	router *gin.Engine
}

func NewServer(st db.Store) Server {
	router := gin.Default()
	server := Server{
		store:  st,
		router: router,
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	router.POST("/users", server.createUser)

	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.listAccounts)

	router.POST("/transfers", server.createTransfer)
	server.router = router

	return server
}

func (s *Server) Start(add string) error {
	return s.router.Run(add)
}

func (s *Server) errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
