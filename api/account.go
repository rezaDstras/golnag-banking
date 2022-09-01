package api

import (
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"github.com/rezaDastrs/banking/db/sqlc"
	"github.com/rezaDastrs/banking/token"
	"net/http"
)

type createAccountRequest struct {
	Currency string `json:"currency" binding:"required,oneof=USD EUR GBP"`
}
func (server *Server)createAccount(ctx *gin.Context)  {
	var req createAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil{
		ctx.JSON(http.StatusBadRequest,errorResponse(err))
		return
	}

	//get authenticated user from payload and cast to payload object
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	arg :=db.CreateAccountParams{
		Owner:    authPayload.Username,
		Balance:  0,
		Currency: req.Currency,
	}

	account , err := server.store.CreateAccount(ctx,arg)
	if err != nil{
		if pqErr , ok := err.(*pq.Error);ok{
			switch pqErr.Code.Name() {
			case "foregin_key_validation","unique_violation":
				ctx.JSON(http.StatusForbidden,errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError,errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK,account)

}

type getAccountRequest struct {
	ID int64 `uri:"id" binding:"required" min:"1"`
}


func (server *Server)getAccount(ctx *gin.Context) {
	var req getAccountRequest
	if err := ctx.ShouldBindUri(&req); err !=nil{
		ctx.JSON(http.StatusBadRequest,errorResponse(err))
		return
	}

	account , err := server.store.GetAccount(ctx,req.ID)
	if err != nil{
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound,errorResponse(err))
		}
		ctx.JSON(http.StatusInternalServerError,errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if authPayload.Username != account.Owner {
		err := errors.New("account doesn't belongs to the authenticated user")
		ctx.JSON(http.StatusUnauthorized,errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK,account)
	return
}

type getAccountsRequest struct {
	Offset int32 `form:"offset" binding:"required,min:1"`
	Limit int32 `form:"limit" binding:"required,min:5,mux=10"`
}
func (server *Server)getAccounts(ctx *gin.Context) {
	var req getAccountsRequest
	if err := ctx.ShouldBindQuery(&req); err !=nil{
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound,errorResponse(err))
		}
		ctx.JSON(http.StatusBadRequest,errorResponse(err))
		return
	}
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	arg := db.ListAccountsParams{
		Owner : authPayload.Username ,
		Offset: req.Offset,
		Limit:  req.Limit,
	}
	accounts , err := server.store.ListAccounts(ctx,arg)
	if err != nil{
		ctx.JSON(http.StatusInternalServerError,errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK,accounts)
	return


}

