package main

import (
	"log"
	"net/http"
	"pkitool/server"
)

func main() {
	router := server.NewRouter()
	log.Fatal(http.ListenAndServe(":8080", router))
}
