// database/database.go
package database

import (
	"database/sql"
	_ "github.com/lib/pq"
)

var db *sql.DB

func InitDB() error {
	var err error
	db, err = sql.Open("postgres", "postgres://arni:arni1375@localhost/db_1?sslmode=disable")
	if err != nil {
		return err
	}
	return db.Ping()
}

func GetDB() *sql.DB {
	return db
}
