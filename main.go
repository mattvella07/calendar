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
	http.Handle("/api/isValidUser", m.MethodChecker(mw.ValidateCookie(http.HandlerFunc(user.IsValid))))
	http.Handle("/api/getEvents", m.MethodChecker(mw.ValidateCookie(http.HandlerFunc(events.GetByDateRange))))
	http.Handle("/api/getEvent/", m.MethodChecker(mw.ValidateCookie(http.HandlerFunc(events.GetByID))))

	http.Handle("/", http.FileServer(http.Dir("./frontend/dist")))

	m.Allowed = []string{"POST"}
	http.Handle("/api/signup", m.MethodChecker(http.HandlerFunc(user.Create)))
	http.Handle("/api/login", m.MethodChecker(http.HandlerFunc(user.Login)))
	http.Handle("/api/logout", m.MethodChecker(http.HandlerFunc(user.Logout)))
	http.Handle("/api/changePassword", m.MethodChecker(mw.ValidateCookie(http.HandlerFunc(user.ChangePassword))))

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func main() {
	conn.InitDB()
	defer conn.DB.Close()

	conn.InitCache()
	defer conn.Cache.Close()

	createServer()
}
