package database

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

// TODO: don't use singleton
var DB *sql.DB
var err error

func InitDB() {
	DB, err = sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/industry_project")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Datebase succesfully opened")
	}
}
