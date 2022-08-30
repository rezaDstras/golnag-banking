package api

import (
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	db "github.com/rezaDastrs/banking/db/sqlc"
	"github.com/rezaDastrs/banking/util"
	"net/http"
	"time"
)

type createUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
	Fullname string `json:"full_name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
}

type userResponse struct {
	Usrname string `json:"usrname"`
	Fullname string `json:"fullname"`
	Email string `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

func (server *Server ) createUser(ctx *gin.Context)  {
	var req createUserRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest,errorResponse(err))
		return
	}

	//create hash password
	hashPassword ,err  := util.HashPassword(req.Password)
	if err != nil{
		ctx.JSON(http.StatusInternalServerError,errorResponse(err))
		return
	}

	arg := db.CreateUserParams{
		Username:     req.Username,
		HashPassword: hashPassword,
		FullName:     req.Fullname,
		Email:        req.Email,
	}

	user , err := server.store.CreateUser(ctx,arg)
	if err != nil {
		if pqErr , ok := err.(*pq.Error);ok{
			switch pqErr.Code.Name() {
			case "unique_validation":
				ctx.JSON(http.StatusForbidden,errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError,errorResponse(err))
		return
	}
	rsp := userResponse{
		Usrname:   user.Username,
		Fullname:  user.FullName,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}
	ctx.JSON(http.StatusOK,rsp)
}
