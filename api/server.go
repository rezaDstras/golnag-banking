package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/rezaDastrs/banking/db/sqlc"
	"github.com/rezaDastrs/banking/token"
	"github.com/rezaDastrs/banking/util"
)

// Server serves Http requests for our services
type Server struct {
	config util.Config
	store *db.Store
	tokenMaker token.Maker
	router *gin.Engine
}


//NewServer creates a new Http server and setup Routing
func NewServer(config util.Config ,store *db.Store) (*Server , error) {

	//jwt
	//tokenMaker , err := token.NewJWToMaker(config.TokenSymmetricToken)

	//paesto
	tokenMaker , err := token.NewPasetoMaker(config.TokenSymmetricToken)

	if err != nil {
		return nil , fmt.Errorf("can not create token maker : %w",err)
	}


	server := &Server{
		config: config,
		store:  store,
		tokenMaker: tokenMaker,
	}



	//validator by gin
	if v , ok :=binding.Validator.Engine().(*validator.Validate);ok{
		v.RegisterValidation("currency",validCurrency)
	}

	//routes
	server.setupRouter()
	return server , nil
}

func (server *Server) setupRouter()  {
	router := gin.Default()

	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))
	authRoutes.POST("/account",server.createAccount)
	authRoutes.GET("/account/:id",server.getAccount)
	authRoutes.GET("/accounts",server.getAccounts)
	authRoutes.POST("/transfer",server.createTransfer)



	router.POST("/user",server.createUser)
	router.POST("/login",server.loginUser)
	router.POST("/token/renew_access",server.renewAccessToken)

	//add routes to routers
	server.router = router
}

// Start runs the http server on a spesefic address
func (server *Server) Start(addr string) error {
	return server.router.Run(addr)
}

//map of errors
func errorResponse(err error) gin.H {
	return gin.H{
		"error": err.Error(),
	}
}
