package main

import (
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/tamboto2000/gotoko-pos/logger"
)

func main() {
	logging, err := logger.NewLogger(logger.ModeStaging)
	if err != nil {
		log.Fatal(err.Error())
	}

	buildRepositories(logging)
	buildServices(logging)
	r := routes(logging)

	logging.Fatal(http.ListenAndServe(":3030", r).Error())
}
