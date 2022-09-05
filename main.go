package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/rezaDastrs/banking/api"
	db "github.com/rezaDastrs/banking/db/sqlc"
	"github.com/rezaDastrs/banking/gapi"
	"github.com/rezaDastrs/banking/pb"
	"github.com/rezaDastrs/banking/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)



func main()  {
	//config
	config , err := util.LoadConfig(".") //same directory as main.go
	if err != nil{
		log.Fatal("cannot load confid:" , err)
	}
	conn, err := sql.Open(config.DBDriver , config.DBSource)
	if err!= nil{
		log.Fatal("cannot connect to db:",err)
	}

	store := db.NewStore(conn)

	//Grpc Server
	runGrpcServer(config , store)

	//Http Server
	//runGinServer(config , store)

}

func runGrpcServer(config util.Config ,store *db.Store)  {
	//create new grpc object server
	server , err := gapi.NewServer(config , store)
	if err != nil {
		log.Fatal("cannot create server:",err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterSimpleBankServer(grpcServer , server)
	//explore easily by client to call avilable rpcs
	reflection.Register(grpcServer)

	listener , err := net.Listen("tcp",config.GrpcServerAddr)
	if err != nil {
		log.Fatal("can not create listener")
	}

	log.Printf("start gRPC server at %s" , listener.Addr().String())

	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("cannot start gRPC server")
	}
}

func runGinServer(config util.Config , store *db.Store){
	server , err := api.NewServer(config , store)
	if err != nil {
		log.Fatal("cannot create server:",err)
	}

	err = server.Start(config.HttpServerAddr)
	if err != nil{
		log.Fatal("cannot start server:",err)
	}
}


