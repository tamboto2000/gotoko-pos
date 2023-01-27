package main

import (
	"flag"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/tamboto2000/gotoko-pos/logger"
)

func main() {
	localEnv := flag.Bool("lenv", false, "if true, local .env file will be used for configuration")

	flag.Parse()

	if *localEnv {
		log.Println("using local .env")
		if err := godotenv.Load(); err != nil {
			log.Fatal(err.Error())
		}
	}

	logging, err := logger.NewLogger(logger.ModeStaging)
	if err != nil {
		log.Fatal(err.Error())
	}

	buildRepositories(logging)
	buildServices(logging)
	r := routes(logging)

	logging.Fatal(http.ListenAndServe(":3030", r).Error())
}
