package api

import (
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
