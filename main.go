package main

import (
	"log"
	"net/http"

	"github.com/mattvella07/calendar-server/api/conn"
	"github.com/mattvella07/calendar-server/api/events"
	mw "github.com/mattvella07/calendar-server/api/middleware"
	"github.com/mattvella07/calendar-server/api/user"
)

func createServer() {
	m := mw.Method{}

	m.Allowed = []string{"GET"}
	http.Handle("/user/list", m.MethodChecker(mw.ValidateJWT(http.HandlerFunc(user.List))))
	http.Handle("/api/getEvents", m.MethodChecker(mw.ValidateJWT(http.HandlerFunc(events.GetByDateRange))))
	http.Handle("/api/getEvent/", m.MethodChecker(mw.ValidateJWT(http.HandlerFunc(events.GetByID))))
	// http.Handle("/", m.MethodChecker(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
	// 	http.ServeFile(rw, r, "./frontend/dist/index.html")
	// })))

	http.Handle("/", http.FileServer(http.Dir("./frontend/dist")))

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
