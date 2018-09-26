package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func InitDatabase() *sql.DB {
	db := db()
	return db
}

func db() *sql.DB {
	err := godotenv.Load()
	if err != nil {
		log.Println(err)
	}
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	_db := os.Getenv("DB")

	db, _ := sql.Open("mysql", user+":"+password+"@tcp(127.0.0.1:3306)/"+_db)
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("Database opened.")
	return db
}
