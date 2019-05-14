package main

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/mattvella07/calendar-server/api/conn"
	"github.com/mattvella07/calendar-server/api/events"
	mw "github.com/mattvella07/calendar-server/api/middleware"
	"github.com/mattvella07/calendar-server/api/user"
)

// Docs contains all info for endpoints
type Docs struct {
	Endpoints []Endpoint `json:"endpoints"`
}

// Endpoint represents an API endpoint
type Endpoint struct {
	URL 				string `json:"url"`
	Method 			string `json:"method"`
	Params 			string `json:"params"`
	ReturnVal 	string `json:"returnVal"`
	Description string `json:"description"`
}

func createServer() {
	m := mw.Method{}

	m.Allowed = []string{"GET"}
	http.Handle("/api", m.MethodChecker(http.HandlerFunc(getDocs)))
	http.Handle("/api/isValidUser", m.MethodChecker(mw.ValidateCookie(http.HandlerFunc(user.IsValid))))
	http.Handle("/api/getEvents", m.MethodChecker(mw.ValidateCookie(http.HandlerFunc(events.GetByDateRange))))
	http.Handle("/api/getEvent/", m.MethodChecker(mw.ValidateCookie(http.HandlerFunc(events.GetByID))))

	http.Handle("/", http.FileServer(http.Dir("./frontend/dist")))

	m.Allowed = []string{"POST"}
	http.Handle("/api/signup", m.MethodChecker(http.HandlerFunc(user.Create)))
	http.Handle("/api/login", m.MethodChecker(http.HandlerFunc(user.Login)))
	http.Handle("/api/logout", m.MethodChecker(http.HandlerFunc(user.Logout)))
	http.Handle("/api/createEvent", m.MethodChecker(mw.ValidateCookie(http.HandlerFunc(events.Create))))

	m.Allowed = []string{"PUT"}
	http.Handle("/api/changePassword", m.MethodChecker(mw.ValidateCookie(http.HandlerFunc(user.ChangePassword))))

	m.Allowed = []string{"DELETE"}
	http.Handle("/api/deleteAccount", m.MethodChecker(mw.ValidateCookie(http.HandlerFunc(user.Delete))))
	http.Handle("/api/deleteEvent/", m.MethodChecker(mw.ValidateCookie(http.HandlerFunc(events.Delete))))

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func main() {
	conn.InitDB()
	defer conn.DB.Close()

	conn.InitCache()
	defer conn.Cache.Close()

	createServer()
}

func getDocs(rw http.ResponseWriter, r *http.Request) {
	jsonFile, err := os.Open("docs.json")
	if err != nil {
		log.Printf("Error opening docs.json: %s\n", err)
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("Error getting documentation"))
		return
	}

	data, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		log.Printf("Error reading from docs.json: %s\n", err)
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("Error getting documentation"))
		return
	}

	d := Docs{}
	err = json.Unmarshal(data, &d)
	if err != nil {
		log.Printf("Error unmarshalling docs.json: %s\n", err)
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("Error getting documentation"))
		return
	}

	tmpl := template.Must(template.ParseFiles("docs.html"))
	tmpl.Execute(rw, d)
}
