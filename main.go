package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/rezaDastrs/banking/api"
	db "github.com/rezaDastrs/banking/db/sqlc"
	"log"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
	serverAddr = "localhost:8080"
)

func main()  {

	conn, err := sql.Open(dbDriver , dbSource)
	if err!= nil{
		log.Fatal("cannot connect to db:",err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(serverAddr)
	if err != nil{
		log.Fatal("cannot start server:",err)
	}
}
