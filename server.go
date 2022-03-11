package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func handleRequests() {
	http.HandleFunc("/", router.getCourse())
}

func SetupDB() {
	var Db *sql.DB
	var err error
	Db, err = sql.Open("mysql", "root:1234@tcp(127.0.0.1:3306)/coursedb")
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Database is Connected!!")
	}

	Db.SetConnMaxLifetime(time.Minute * 3)
	Db.SetMaxOpenConns(10)
	Db.SetMaxIdleConns(10)
}

func main() {
	SetupDB()
	log.Fatal(http.ListenAndServe(":80", nil))
}
