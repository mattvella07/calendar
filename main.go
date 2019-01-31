package main

import (
	"log"
	"net/http"

	"github.com/mattvella07/calendar-server/conn"
	"github.com/mattvella07/calendar-server/user"
)

func createServer() {
	http.Handle("/user/create", http.HandlerFunc(user.Create))

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func main() {
	conn.InitDB()
	defer conn.DB.Close()

	createServer()
}
