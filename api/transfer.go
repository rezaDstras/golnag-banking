package api

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rezaDastrs/banking/db/sqlc"
	"net/http"
)

type createTransferRequest struct {
	FromAccountID int64 `json:"from_account" binding:"required" min:"1"`
	ToAccountID   int64 `json:"to_account" binding:"required" min:"1"`
	Amount      int64 `json:"amount" binding:"required" gt:"0"`
	Currency   string `json:"currency" binding:"required,currency"`
}

func (server *Server)createTransfer(ctx *gin.Context) {
	var req createTransferRequest
	if err := ctx.ShouldBindJSON(&req); err !=nil{
		ctx.JSON(http.StatusBadRequest,errorResponse(err))
		return
	}

	if !server.validateCurrency(ctx,req.FromAccountID,req.Currency) {
		return
	}

	if server.validateCurrency(ctx,req.ToAccountID,req.Currency) {
		return
	}
	
	arg := db.CreateTransferParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}

	result , err := server.store.CreateTransfer(ctx,arg)
	if err !=nil {
		ctx.JSON(http.StatusInternalServerError,errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK,result)
	return
}

func (server *Server) validateCurrency(ctx *gin.Context , accountID int64 , currency string ) bool {
	account , err := server.store.GetAccount(ctx,accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound,errorResponse(err))
			return false
		}
		ctx.JSON(http.StatusInternalServerError,errorResponse(err))
	}

	if account.Currency != currency {
		err:=fmt.Errorf("account [%d] currency mismatch %s vs %s ",accountID,account.Currency,currency)
		ctx.JSON(http.StatusBadRequest,errorResponse(err))
		return false
	}

	return true

}