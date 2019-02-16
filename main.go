package main

import (
	"log"
	"net/http"

	"github.com/mattvella07/calendar-server/conn"
	"github.com/mattvella07/calendar-server/events"
	mw "github.com/mattvella07/calendar-server/middleware"
	"github.com/mattvella07/calendar-server/user"
)

func createServer() {
	m := mw.Method{}

	m.Allowed = []string{"GET"}
	http.Handle("/user/list", m.MethodChecker(mw.ValidateJWT(http.HandlerFunc(user.List))))
	http.Handle("/api/getEvents", m.MethodChecker(mw.ValidateJWT(http.HandlerFunc(events.Get))))

	m.Allowed = []string{"POST"}
	http.Handle("/api/signup", m.MethodChecker(http.HandlerFunc(user.Create)))
	http.Handle("/api/login", m.MethodChecker(http.HandlerFunc(user.Login)))

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func main() {
	conn.InitDB()
	defer conn.DB.Close()

	createServer()
}
