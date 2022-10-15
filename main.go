package main

import (
	"database/sql"
	"log"
	"simplebank/api"
	db "simplebank/db/sqlc"

	_ "github.com/lib/pq"
)

const (
	dbdriver      = "postgres"
	dbsource      = "postgresql://root:password@localhost:8080/simple_bank?sslmode=disable"
	serverAddress = ":7070"
)

func main() {
	conn, err := sql.Open(dbdriver, dbsource)
	if err != nil {
		log.Fatalln(err)
	}
	store := db.NewStore(conn)
	server := api.NewServer(&store)

	err = server.Start(serverAddress)
	if err != nil {
		log.Fatalln("cannot start server", err)
	}
}
