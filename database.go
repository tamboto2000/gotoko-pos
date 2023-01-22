package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

const (
	mysqlHostEnv     = "MYSQL_HOST"
	mysqlPortEnv     = "MYSQL_PORT"
	mysqlUserEnv     = "MYSQL_USER"
	mysqlPasswordEnv = "MYSQL_PASSWORD"
	mysqlDbNameEnv   = "MYSQL_DBNAME"
)

func buildDatabase() (*sql.DB, error) {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err.Error())
	}

	host := os.Getenv(mysqlHostEnv)
	port := os.Getenv(mysqlPortEnv)
	user := os.Getenv(mysqlUserEnv)
	pass := os.Getenv(mysqlPasswordEnv)
	dbname := os.Getenv(mysqlDbNameEnv)

	connStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", user, pass, host, port, dbname)

	return sql.Open("mysql", connStr)
}
