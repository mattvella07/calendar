package conn

import (
	"database/sql"
	"fmt"
	"log"
)

// DB is the reference to the database
var DB *sql.DB

//InitDB initializes the database connection
func InitDB() {
	var err error
	DB, err = sql.Open("postgres", fmt.Sprintf("host=localhost port=%d user=postgres password=postgres dbname=calendar sslmode=disable", 5432))
	if err != nil {
		log.Fatalf("Error connecting to DB: %s", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatalf("Error communicating with DB: %s", err)
	}
}
