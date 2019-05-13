package events

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	startDate := r.URL.Query().Get("startDate")
	endDate := r.URL.Query().Get("endDate")

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

	if strings.TrimSpace(eventID) == "" {
		log.Printf("Missing event id\n")
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Event ID must be provided"))
		return
	}

	query := `SELECT * FROM events
		WHERE owner_id = $1 
		AND id = $2`

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

// Create creates a new event
func Create(rw http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("userid")

	// Get event data from body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading request body: %s\n", err)
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Error reading request body"))
		return
	}

	e := Event{}
	err = json.Unmarshal(body, &e)
	if err != nil {
		log.Printf("Error reading request body: %s\n", err)
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Error reading request body"))
		return
	}

	// Validate event data 
	if strings.TrimSpace(e.Title) == "" || strings.TrimSpace(e.StartTime) == "" || strings.TrimSpace(e.EndTime) == "" {
		log.Printf("Invalid params: title, startDate or endDate is missing\n")
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Missing title, startDate or endDate"))
		return
	}

	if _, err := time.Parse(time.RFC3339, e.StartTime); err != nil {
		log.Printf("Invalid value for startDate\n")
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Invalid startDate"))
		return
	}

	if _, err := time.Parse(time.RFC3339, e.EndTime); err != nil {
		log.Printf("Invalid value for endDate\n")
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Invalid endDate"))
		return
	}

	// Insert into DB
	eventID := -1
	query := `INSERT INTO events
		(title, "start_time", "end_time", location, notes, "owner_id")
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id`

	err = conn.DB.QueryRow(query, e.Title, e.StartTime, e.EndTime, e.Location, e.Notes, userID).Scan(&eventID)
	if err != nil {
		log.Printf("Error creating event: %s\n", err)
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("Error creating event"))
		return
	}

	log.Printf("Event %s created\n", e.Title)
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte(fmt.Sprintf("Event %s created", e.Title)))
}

// Delete deletes the specified event
func Delete(rw http.ResponseWriter, r *http.Request) {
	// Get userID and eventID
	userID := r.Header.Get("userid")
	urlParams := strings.Replace(r.URL.String(), "/api/deleteEvent/", "", 1)
	eventID := strings.Split(urlParams, "/")[0]

	// eventID must not be empty
	if strings.TrimSpace(eventID) == "" {
		log.Printf("Missing event id to delete\n")
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Event ID must be provided"))
		return
	}

	// Delete from DB
	query := `DELETE
		FROM events
		WHERE id = $1
		AND owner_id = $2`

	res, err := conn.DB.Exec(query, eventID, userID)
	if err != nil {
		log.Printf("Error deleting event: %s\n", err)
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("Error deleting event"))
		return
	}

	if count, err := res.RowsAffected(); err != nil || count == 0 {
		log.Printf("Error deleting event, count: %d, err: %s\n", count, err)
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("Error deleting event"))
		return
	}

	log.Println("Event deleted")
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("Event deleted"))
}
