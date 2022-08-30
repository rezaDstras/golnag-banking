package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/rezaDastrs/banking/db/sqlc"
)

// Server serves Http requests for our services
type Server struct {
	store *db.Store
	router *gin.Engine
}

//NewServer creates a new Http server and setup Routing
func NewServer(store *db.Store) *Server  {
	server := &Server{
		store:  store,
	}

	router := gin.Default()

	//validator by gin
	if v , ok :=binding.Validator.Engine().(*validator.Validate);ok{
		v.RegisterValidation("currency",validCurrency)
	}

	//routes
	router.POST("/account",server.createAccount)
	router.GET("/account/:id",server.getAccount)
	router.GET("/accounts",server.getAccounts)
	router.POST("/transfer",server.createTransfer)

	//add routes to routers
	server.router = router
	return server
}

//Starts runs the http server on a spesefic address
func (server *Server) Start(addr string) error {
	return server.router.Run(addr)
}

//map of errors
func errorResponse(err error) gin.H {
	return gin.H{
		"error": err.Error(),
	}
}
