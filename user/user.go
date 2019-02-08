package user

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	jwt "github.com/dgrijalva/jwt-go"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"

	"github.com/mattvella07/calendar-server/conn"
)

// User contains info about the user
type User struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

// Create creates a new user
func Create(rw http.ResponseWriter, r *http.Request) {
	username, password, ok := r.BasicAuth()
	if !ok {
		log.Println("Username and/or password missing")
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Username and/or password missing"))
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading request body: %s\n", err)
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Error reading request body"))
		return
	}

	u := User{}
	err = json.Unmarshal(body, &u)
	if err != nil {
		log.Printf("Error reading request body: %s\n", err)
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Error reading request body"))
		return
	}

	existingUser := ""
	err = conn.DB.QueryRow(`SELECT username FROM users WHERE username = $1`, username).Scan(&existingUser)

	switch {
	case err == sql.ErrNoRows:
		// User doesn't already exist, add user
		passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), 14)
		if err != nil {
			log.Printf("Error creating user: %s\n", err)
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte("Error creating user"))
			return
		}

		_, err = conn.DB.Exec(`INSERT INTO users (username, password, "first_name", "last_name") VALUES ($1, $2, $3, $4)`, username, passwordHash, u.FirstName, u.LastName)
		if err != nil {
			log.Printf("Error creating user: %s\n", err)
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte("Error creating user"))
			return
		}

		log.Printf("User %s created", username)
		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte(fmt.Sprintf("User %s created", username)))
	case err != nil:
		// Other error
		log.Printf("Error: %s\n", err)
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("Error creating user"))
	default:
		// User already exists
		log.Printf("User %s already exists\n", username)
		rw.WriteHeader(http.StatusForbidden)
		rw.Write([]byte(fmt.Sprintf("User %s already exists", username)))
	}
}

// Login validates the username and password and logs the user in
func Login(rw http.ResponseWriter, r *http.Request) {
	username, password, ok := r.BasicAuth()
	if !ok {
		log.Println("Username and/or password missing")
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Username and/or password missing"))
		return
	}

	existingUser := User{}
	err := conn.DB.QueryRow(`SELECT username, password FROM users WHERE username = $1 LIMIT 1`, username).Scan(&existingUser.Username, &existingUser.Password)
	if err != nil {
		log.Printf("DB error: %s\n", err)
		rw.WriteHeader(http.StatusUnauthorized)
		rw.Write([]byte("Incorrect Username and/or password"))
		return
	}

	// User exists, valdiate password
	err = bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(password))
	if err != nil {
		log.Printf("Incorrect password: %s\n", err)
		rw.WriteHeader(http.StatusUnauthorized)
		rw.Write([]byte("Incorrect Username and/or password"))
		return
	}

	token := jwt.New(jwt.SigningMethodHS256)
	signingKey := os.Getenv("SIGNING_KEY")
	claims := token.Claims.(jwt.MapClaims)

	claims["name"] = "MV"

	tokenStr, _ := token.SignedString([]byte(signingKey))

	fmt.Println(tokenStr)

	log.Printf("User %s logged in\n", username)
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("Success"))
}

// List will list all users
func List(rw http.ResponseWriter, r *http.Request) {
	rows, err := conn.DB.Query(`SELECT username, password, first_name, last_name FROM users`)
	if err != nil {
		log.Printf("DB error: %s\n", err)
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Unable to communicate with database"))
		return
	}
	defer rows.Close()

	allUsers := []User{}
	for rows.Next() {
		u := User{}
		err = rows.Scan(&u.Username, &u.Password, &u.FirstName, &u.LastName)
		if err != nil {
			log.Printf("DB error: %s\n", err)
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write([]byte("Error reading from database"))
			return
		}
		allUsers = append(allUsers, u)
	}

	rw.WriteHeader(http.StatusOK)
	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(allUsers)
}
