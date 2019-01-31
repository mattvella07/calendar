package user

import (
	"database/sql"
	"fmt"
	"net/http"

	_ "github.com/lib/pq"

	"github.com/mattvella07/calendar-server/conn"
)

// Create creates a new user
func Create(rw http.ResponseWriter, r *http.Request) {
	usr := ""
	err := conn.DB.QueryRow(`SELECT username FROM users LIMIT 1`).Scan(&usr)

	switch {
	case err == sql.ErrNoRows:
		// User doesn't exist
		fmt.Println("No Rows")
	case err != nil:
		// Other error
		fmt.Printf("Err: %s", err)
	default:
		// User already exists
		fmt.Println("User ", usr)
	}
}
