package gapi

import (
	"fmt"
	db "github.com/rezaDastrs/banking/db/sqlc"
	"github.com/rezaDastrs/banking/pb"
	"github.com/rezaDastrs/banking/token"
	"github.com/rezaDastrs/banking/util"
)

// Server serves Grpc requests for our services
type Server struct {
	pb.UnimplementedSimpleBankServer
	config util.Config
	store *db.Store
	tokenMaker token.Maker
}


//NewServer creates a new Grpc server and setup Routing
func NewServer(config util.Config ,store *db.Store) (*Server , error) {


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
 
	return server , nil
}