package api

import (
	"database/sql"
	"net/http"
	db "simplebank/db/sqlc"
	"simplebank/util"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type createUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
	Fullname string `json:"fullname" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

type userResponse struct {
	Username          string    `json:"username"`
	Fullname          string    `json:"fullname"`
	Email             string    `json:"email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

func newUserResponse(user db.User) userResponse {
	return userResponse{
		Username:          user.Username,
		Fullname:          user.Fullname,
		Email:             user.Email,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
	}
}

func (server *Server) createUser(c *gin.Context) {
	var req createUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, server.errorResponse(err))
		return
	}

	hashedPassord, err := util.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, server.errorResponse(err))
		return
	}

	arg := db.CreateUserParams{
		Username:       req.Username,
		HashedPassword: hashedPassord,
		Email:          req.Email,
		Fullname:       req.Fullname,
	}

	user, err := server.store.CreateUser(c, arg)
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok {
			switch pgErr.Code.Name() {
			case "unique_violation":
				c.JSON(http.StatusForbidden, server.errorResponse(err))
				return
			}
		}
		c.JSON(http.StatusInternalServerError, server.errorResponse(err))
		return
	}

	res := newUserResponse(user)

	c.JSON(http.StatusOK, res)
}

type loginUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
}

type loginUserResponse struct {
	AccessToken string       `json:"access_token"`
	User        userResponse `json:"user"`
}

func (server *Server) loginUser(c *gin.Context) {
	var req loginUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, server.errorResponse(err))
		return
	}

	user, err := server.store.GetUser(c, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, server.errorResponse(err))
			return
		}
		c.JSON(http.StatusInternalServerError, server.errorResponse(err))
		return
	}

	err = util.CompareHashAndPassword(user.HashedPassword, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, server.errorResponse(err))
		return
	}

	accessToken, err := server.tokenMaker.CreateToken(
		req.Username,
		server.config.TokenDuration,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, server.errorResponse(err))
		return
	}

	res := loginUserResponse{
		AccessToken: accessToken,
		User:        newUserResponse(user),
	}

	c.JSON(http.StatusOK, res)
}
