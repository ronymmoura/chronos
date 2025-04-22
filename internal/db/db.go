package db

import (
	"fmt"
	"log"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/jmoiron/sqlx"
)

type DB struct {
	Conn *sqlx.DB
}

func Connect(connectionString string) *DB {
	conn, err := sqlx.Connect("mssql", connectionString)
	if err != nil {
		log.Fatalf("cannot connect to db: %v", err)
	}

	fmt.Println("db connected")

	return &DB{Conn: conn}
}
