package main

import (
	"forum/DB"
	"forum/handlers"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	srvr := http.Server{
		Addr:    ":443",
		Handler: handlers.Routes(),
	}

	DB.InitDB()
	log.Println("starting server on https://localhost/")
	err := srvr.ListenAndServeTLS("./cert/cert.pem", "./cert/key.pem")
	if err != nil {
		log.Fatalf("error starting server:%v", err)
	}
}
