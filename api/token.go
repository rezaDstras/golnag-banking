package api

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type renewAccessTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type renewAccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	AccessTokenExpireAt time.Time `json:"access_token_expire_at"`
}

func (server *Server) renewAccessToken (ctx *gin.Context)  {
	var req renewAccessTokenRequest
	if err := ctx.ShouldBindJSON(&req);err !=nil{
		ctx.JSON(http.StatusBadRequest,errorResponse(err))
	}

	refreshPayload ,err := server.tokenMaker.VerifyToken(req.RefreshToken)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized,errorResponse(err))
		return
	}
	session, err := server.store.GetSession(ctx , refreshPayload.ID)
	if err != nil{
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound,errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError,errorResponse(err))
		return
	}

	if session.IsBlocked {
		err := fmt.Errorf("seesion is blocked")
		ctx.JSON(http.StatusUnauthorized,errorResponse(err))
		return
	}

	if session.Username != refreshPayload.Username {
		err := fmt.Errorf("incorrect session user")
		ctx.JSON(http.StatusUnauthorized,errorResponse(err))
		return
	}

	if session.RefreshToken !=req.RefreshToken {
		err := fmt.Errorf("missmatch session token")
		ctx.JSON(http.StatusUnauthorized,errorResponse(err))
		return
	}

	if time.Now().After(session.ExpireAt) {
		err := fmt.Errorf("expired session")
		ctx.JSON(http.StatusUnauthorized,errorResponse(err))
		return
	}

	//create access token
	accessToken , accessPayload , err := server.tokenMaker.CreateToken(
		session.Username,
		server.config.AccessTokenDuration,
		)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError,errorResponse(err))
	}

	rep:= renewAccessTokenResponse{
		AccessToken:         accessToken,
		AccessTokenExpireAt: accessPayload.ExpireAt,
	}

	ctx.JSON(http.StatusOK,rep)

}
