package database

import (
    "database/sql"
	"fmt"
    "os"
    _ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func ConnectToDatabase() error {
    var err error
    db_pass := os.Getenv("DB_PASS")
    db_user := os.Getenv("DB_USER")
    db_name := os.Getenv("DB_NAME")
    db, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(localhost:3306)/%s", db_user, db_pass, db_name))
    if err != nil {
        panic(err.Error())
    }

    err = db.Ping()
    if err != nil {
        panic(err.Error())
    }
    fmt.Println("Successfully connected to db")
    return nil
}

func GetDB() *sql.DB {
    return db
}
