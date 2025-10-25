package db

import (
	"database/sql"
	"fmt"
	"log"
	_ "github.com/lib/pq"
)


var DB *sql.DB

func InitDB(dataSourceName string) {
	var err error 
	DB, err = sql.Open("postgres", dataSourceName)
	if err != nil {
		log.Fatal("Error connecting to the database: ", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal("Error pinging the database: ", err)
		
	}

	fmt.Println("Successfully connected to the database")
}