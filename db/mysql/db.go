package mysql

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type Cliente struct {
	Id     int    `json:"id"`
	Nome   string `json:"nome"`
	Limite int    `json:"limite"`
}

func Init() *sql.DB {
	fmt.Println("Connecting to MySQL...")

	dbHost := os.Getenv("DB_HOSTNAME")
	if dbHost == "" {
		dbHost = "127.0.0.1"
	}

	connection := fmt.Sprintf("rinha:secretpass@tcp(%s:3306)/rinhabackend?parseTime=true", dbHost)
	db, err := sql.Open("mysql", connection)

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(8)

	if err != nil {
		fmt.Println("Error connecting to MySQL:", err)
		panic(err.Error())
	}
	fmt.Println("Connected to MySQL!")

	return db
}
