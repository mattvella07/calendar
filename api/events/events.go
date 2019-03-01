package events

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/mattvella07/calendar-server/api/conn"
)

// Event contains info about an event
type Event struct {
	EventID   int    `json:"id"`
	OwnerID   int    `json:"owner_id"`
	Title     string `json:"title"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	Location  string `json:"location"`
	Notes     string `json:"notes"`
}

// GetByDateRange returns all events for a specific user and date range
func GetByDateRange(rw http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("userid")
	startDate := r.Header.Get("startDate")
	endDate := r.Header.Get("endDate")

	// Validate startDate and endDate
	if startDate == "" || endDate == "" {
		log.Printf("Invalid params, startDate or endDate is missing\n")
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Missing startDate or endDate"))
		return
	}

	if _, err := time.Parse(time.RFC3339, startDate); err != nil {
		log.Printf("Invalid value for startDate\n")
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Invalid startDate"))
		return
	}

	if _, err := time.Parse(time.RFC3339, endDate); err != nil {
		log.Printf("Invalid value for endDate\n")
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Invalid endDate"))
		return
	}

	query := `SELECT * FROM events
		WHERE owner_id = $1 AND
		(start_time BETWEEN $2 AND $3 OR
		end_time BETWEEN $2 AND $3)`

	rows, err := conn.DB.Query(query, userID, startDate, endDate)
	if err != nil {
		log.Printf("DB error: %s\n", err)
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Unable to communicate with database"))
		return
	}
	defer rows.Close()

	events := []Event{}
	for rows.Next() {
		e := Event{}
		loc := sql.NullString{}
		notes := sql.NullString{}

		err = rows.Scan(&e.EventID, &e.Title, &e.StartTime, &e.EndTime, &loc, &notes, &e.OwnerID)
		if err != nil {
			log.Printf("DB error: %s\n", err)
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write([]byte("Error reading from database"))
			return
		}

		e.Location = loc.String
		e.Notes = notes.String

		events = append(events, e)
	}

	rw.WriteHeader(http.StatusOK)
	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(events)
}

// GetByID returns the event by ID as long as the user has access to it
func GetByID(rw http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("userid")

	urlParams := strings.Replace(r.URL.String(), "/api/getEvent/", "", 1)
	eventID := strings.Split(urlParams, "/")[0]

	query := `SELECT * FROM events
		WHERE owner_id = $1 AND
		id = $2`

	e := Event{}
	loc := sql.NullString{}
	notes := sql.NullString{}

	err := conn.DB.QueryRow(query, userID, eventID).Scan(&e.EventID, &e.Title, &e.StartTime, &e.EndTime, &loc, &notes, &e.OwnerID)
	if err != nil {
		log.Printf("DB error: %s\n", err)
		rw.WriteHeader(http.StatusNoContent)
		rw.Write([]byte("Event not found"))
		return
	}

	e.Location = loc.String
	e.Notes = notes.String

	rw.WriteHeader(http.StatusOK)
	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(e)
}
