package conn

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gomodule/redigo/redis"
)

// DB is the reference to the database
var DB *sql.DB

// Cache is the reference to the redis cache
var Cache redis.Conn

// InitDB initializes the database connection
func InitDB() {
	var err error
	DB, err = sql.Open("postgres", fmt.Sprintf("host=db port=%d user=%s password=%s dbname=%s sslmode=disable", 5432, os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_DB_NAME")))
	if err != nil {
		log.Fatalf("Unable to connect to DB: %s", err)
	}

	retries := 15
	for retries >= 0 {
		err = DB.Ping()
		if err != nil {
			if retries == 0 {
				log.Fatalf("Unable to communicate with DB: %s", err)
			}

			log.Print("Error communicating with DB, trying again...")
			time.Sleep(1 * time.Second)

			retries--
		} else {
			retries = -1
		}
	}
}

// InitCache initializes the redis connection
func InitCache() {
	var err error
	Cache, err = redis.DialURL("redis://cache")
	if err != nil {
		log.Fatalf("Unable to connect to Redis: %s", err)
	}

	retries := 15
	for retries >= 0 {
		err = Cache.Send("PING")
		if err != nil {
			if retries == 0 {
				log.Fatalf("Unable to communicate with Redis: %s", err)
			}

			log.Print("Error communicating with Redis, trying again...")
			time.Sleep(1 * time.Second)

			retries--
		} else {
			retries = -1
		}
	}
}
