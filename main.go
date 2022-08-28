package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/rezaDastrs/banking/api"
	db "github.com/rezaDastrs/banking/db/sqlc"
	"github.com/rezaDastrs/banking/util"
	"log"
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
	server := api.NewServer(store)

	err = server.Start(config.ServerAddr)
	if err != nil{
		log.Fatal("cannot start server:",err)
	}
}
