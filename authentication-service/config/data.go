package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func ConnectDB() *sql.DB {
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Panic("Error loading .env file")
	// }

	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	user := os.Getenv("USER_NAME")
	password := os.Getenv("PASSWORD")
	dbname := os.Getenv("DB_NAME")
	psqlInfo := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	fmt.Print()
	if err != nil {
		log.Panic(err)
	}
	// defer db.Close()
	err = db.Ping()
	if err != nil {
		log.Panic(err)
	}

	return db
}
