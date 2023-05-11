package main

import (
	"log"
	"net/http"
	database "social-network/back-end/database"
)

type application struct {
}

func main() {
	var app application
	database.ConnectDB()
	defer database.CloseDB()

	log.Println("server started at: http://localhost:8080/")
	log.Fatal(http.ListenAndServe(":8080", app.server()))
}
