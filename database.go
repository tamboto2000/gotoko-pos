package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/tamboto2000/gotoko-pos/migrations"
)

const (
	mysqlHostEnv     = "MYSQL_HOST"
	mysqlPortEnv     = "MYSQL_PORT"
	mysqlUserEnv     = "MYSQL_USER"
	mysqlPasswordEnv = "MYSQL_PASSWORD"
	mysqlDbNameEnv   = "MYSQL_DBNAME"
)

func buildDatabase() (*sql.DB, error) {
	host := os.Getenv(mysqlHostEnv)
	port := os.Getenv(mysqlPortEnv)
	user := os.Getenv(mysqlUserEnv)
	pass := os.Getenv(mysqlPasswordEnv)
	dbname := os.Getenv(mysqlDbNameEnv)

	connStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&multiStatements=true", user, pass, host, port, dbname)
	log.Println("connStr: ", connStr)

	db, err := sql.Open("mysql", connStr)
	if err != nil {
		return nil, err
	}

	if err := migrations.Migrate(db, dbname); err != nil {
		return nil, err
	}

	return db, nil
}
