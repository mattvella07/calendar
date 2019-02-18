package events

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/mattvella07/calendar-server/conn"
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

// Get returns all events for a specific user and date range
func Get(rw http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("userid")

	// Need to filter by a date range (will get from Params)
	startDate := r.Header.Get("startDate")
	endDate := r.Header.Get("endDate")

	// time.Parse() -- check if startDate and endDate are valid times
	if startDate == "" || endDate == "" {
		log.Printf("Invalid params, startDate or endDate is missing:\n")
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Missing startDate or endDate"))
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
