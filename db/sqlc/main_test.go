package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDB *sql.DB

const (
	dbdriver = "postgres"
	dbsource = "postgresql://root:password@localhost:8080/simple_bank?sslmode=disable"
)

func TestMain(m *testing.M) {
	var err error
	testDB, err = sql.Open(dbdriver, dbsource)
	if err != nil {
		log.Fatalln(err)
	}
	testQueries = New(testDB)

	os.Exit(m.Run())
}
