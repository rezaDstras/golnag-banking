package api

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/rezaDastrs/banking/db/sqlc"
	"net/http"
)

type createAccountRequest struct {
	Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,oneof=USD EUR GBP"`
}
func (server *Server)createAccount(ctx *gin.Context)  {
	var req createAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil{
		ctx.JSON(http.StatusBadRequest,errorResponse(err))
		return
	}

	arg :=db.CreateAccountParams{
		Owner:    req.Owner,
		Balance:  0,
		Currency: req.Currency,
	}

	account , err := server.store.CreateAccount(ctx,arg)
	if err != nil{
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
	ctx.JSON(http.StatusOK,account)
	return
}

type getAccountsRequest struct {
	Offset int32 `form:"offset" min:"1"`
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
	accounts , err := server.store.ListAccounts(ctx,req.Offset)
	if err != nil{
		ctx.JSON(http.StatusInternalServerError,errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK,accounts)
	return


}

