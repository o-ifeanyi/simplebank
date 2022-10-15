package api

import (
	db "simplebank/db/sqlc"

	"github.com/gin-gonic/gin"
)

type Server struct {
	store  *db.Store
	router *gin.Engine
}

func NewServer(st *db.Store) Server {
	router := gin.Default()
	server := Server{
		store:  st,
		router: router,
	}

	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.listAccounts)

	server.router = router

	return server
}

func (s *Server) Start(add string) error {
	return s.router.Run(add)
}

func (s *Server) errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
