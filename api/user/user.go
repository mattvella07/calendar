package user

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"

	"github.com/mattvella07/calendar-server/api/conn"
	uuid "github.com/satori/go.uuid"
)

// User contains info about the user
type User struct {
	UserID    int    `josn:"id"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

// IsValid checks if the user is a valid user
func IsValid(rw http.ResponseWriter, r *http.Request) {
	// If ValidateCookie middleware passed then the user is valid
	log.Println("Valid user")
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("Valid user"))
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

	if strings.TrimSpace(username) == "" {
		log.Println("Empty username")
		rw.WriteHeader(http.StatusUnauthorized)
		rw.Write([]byte("Username must not be empty"))
		return
	}

	if strings.TrimSpace(password) == "" {
		log.Println("Empty password")
		rw.WriteHeader(http.StatusUnauthorized)
		rw.Write([]byte("Password must not be empty"))
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

	currUserName := ""
	query := `SELECT username
		FROM users
		WHERE username = $1`

	err = conn.DB.QueryRow(query, username).Scan(&currUserName)

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

		if strings.TrimSpace(u.FirstName) == "" {
			log.Println("Empty first name")
			rw.WriteHeader(http.StatusUnauthorized)
			rw.Write([]byte("First name must not be empty"))
			return
		}

		if strings.TrimSpace(u.LastName) == "" {
			log.Println("Empty last name")
			rw.WriteHeader(http.StatusUnauthorized)
			rw.Write([]byte("Last name must not be empty"))
			return
		}

		userID := -1
		query = `INSERT INTO users
			(username, password, "first_name", "last_name")
			VALUES ($1, $2, $3, $4)
			RETURNING id`

		err = conn.DB.QueryRow(query, username, passwordHash, u.FirstName, u.LastName).Scan(&userID)
		if err != nil {
			log.Printf("Error creating user: %s\n", err)
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte("Error creating user"))
			return
		}

		// token := generateToken(userID, username)

		cookie, err := generateCookie(string(userID))
		if err != nil {
			log.Printf("Error generating cookie: %s\n", err)
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte("Error creating user"))
			return
		}
		http.SetCookie(rw, &cookie)

		log.Printf("User %s created\n", username)
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
	fmt.Println("Username: ", username)
	fmt.Println("Password: ", password)
	if !ok {
		log.Println("Username and/or password missing")
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Username and/or password missing"))
		return
	}

	existingUser := User{}
	query := `SELECT id, username, password
		FROM users
		WHERE username = $1
		LIMIT 1`

	err := conn.DB.QueryRow(query, username).Scan(&existingUser.UserID, &existingUser.Username, &existingUser.Password)
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

	// token := generateToken(existingUser.UserID, existingUser.Username)

	cookie, err := generateCookie(string(existingUser.UserID))
	if err != nil {
		log.Printf("Error generating cookie: %s\n", err)
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("Error logging in"))
		return
	}
	http.SetCookie(rw, &cookie)

	log.Printf("User %s logged in\n", username)
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("Logged in"))
}

// Logout logs the user out
func Logout(rw http.ResponseWriter, r *http.Request) {
	// Remove current token
	sessionToken := r.Header.Get("sessionToken")

	log.Printf("Logout sessionToken: %s", sessionToken)

	_, err := conn.Cache.Do("DEL", sessionToken)
	if err != nil {
		log.Printf("Error deleting session token from cache: %s\n", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Send back empty cookie
	http.SetCookie(rw, &http.Cookie{
		Name:   "session_token",
		Value:  "",
		MaxAge: 0,
	})

	log.Println("User logged out")
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("Logged out"))
}

// ChangePassword changes the user's password
func ChangePassword(rw http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(r.Header.Get("userid"))
	if err != nil {
		log.Printf("Invalid User ID: %s\n", err)
		rw.WriteHeader(http.StatusUnauthorized)
		rw.Write([]byte("Incorrect Username and/or password"))
		return
	}
	currPassword := r.Header.Get("currentPassword")
	newPassword := r.Header.Get("newPassword")
	sessionToken := r.Header.Get("sessionToken")

	// Check if user exists in DB
	existingUser := User{}
	query := `SELECT password
		FROM users
		WHERE id = $1
		LIMIT 1`

	err = conn.DB.QueryRow(query, userID).Scan(&existingUser.Password)
	if err != nil {
		log.Printf("DB error: %s\n", err)
		rw.WriteHeader(http.StatusUnauthorized)
		rw.Write([]byte("Incorrect Username and/or password"))
		return
	}

	// User exists, valdiate current password
	err = bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(currPassword))
	if err != nil {
		log.Printf("Incorrect password: %s\n", err)
		rw.WriteHeader(http.StatusUnauthorized)
		rw.Write([]byte("Incorrect Username and/or password"))
		return
	}

	if strings.TrimSpace(newPassword) == "" {
		log.Println("Empty new password")
		rw.WriteHeader(http.StatusUnauthorized)
		rw.Write([]byte("New password must not be empty"))
		return
	}

	// Update users table with new password
	newPasswordHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), 14)
	if err != nil {
		log.Printf("Error creating hash of user's new password: %s\n", err)
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("Error updating password"))
		return
	}

	query = `UPDATE users
		SET password = $1
		WHERE id = $2`

	res, err := conn.DB.Exec(query, newPasswordHash, userID)
	if err != nil {
		log.Printf("Error updating user's password: %s\n", err)
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("Error updating password"))
		return
	}

	if count, err := res.RowsAffected(); err != nil || count == 0 {
		log.Printf("Error updating user's password, count: %d, err: %s\n", count, err)
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("Error updating password"))
		return
	}

	// Remove current token
	_, err = conn.Cache.Do("DEL", sessionToken)
	if err != nil {
		log.Printf("Error deleting session token from cache: %s\n", err)
	}

	// Delete all tokens for the current user
	deleteAllTokensForUser(userID)

	// Send back new cookie
	cookie, err := generateCookie(string(userID))
	if err != nil {
		log.Printf("Error generating cookie: %s\n", err)

		// If problem generating a new cookie, send back an empty one
		http.SetCookie(rw, &http.Cookie{
			Name:   "session_token",
			Value:  "",
			MaxAge: 0,
		})

		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte("Password updated"))
		return
	}

	http.SetCookie(rw, &cookie)

	log.Println("Password updated")
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("Password updated"))
}

// Delete deletes the current user
func Delete(rw http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(r.Header.Get("userid"))
	if err != nil {
		log.Printf("Invalid User ID: %s\n", err)
		rw.WriteHeader(http.StatusUnauthorized)
		rw.Write([]byte("Incorrect Username and/or password"))
		return
	}
	sessionToken := r.Header.Get("sessionToken")

	// Remove user from users table
	query := `DELETE
		FROM users
		WHERE id = $1`

	res, err := conn.DB.Exec(query, userID)
	if err != nil {
		log.Printf("Error deleting user: %s\n", err)
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("Error deleting user"))
		return
	}

	if count, err := res.RowsAffected(); err != nil || count == 0 {
		log.Printf("Error deleting user, count: %d, err: %s\n", count, err)
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("Error deleting password"))
		return
	}

	// Remove all events owned by that user
	query = `DELETE
		FROM events
		WHERE owner_id = $1`

	_, err = conn.DB.Exec(query, userID)
	if err != nil {
		log.Printf("Error deleting user's events: %s\n", err)
	}

	// Remove current token
	_, err = conn.Cache.Do("DEL", sessionToken)
	if err != nil {
		log.Printf("Error deleting session token from cache: %s\n", err)
	}

	// Remove all existing tokens for that user
	deleteAllTokensForUser(userID)

	// Send back empty cookie
	http.SetCookie(rw, &http.Cookie{
		Name:   "session_token",
		Value:  "",
		MaxAge: 0,
	})

	log.Println("User deleted")
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("User deleted"))
}

func generateCookie(userID string) (http.Cookie, error) {
	sessionToken := uuid.NewV4().String()

	_, err := conn.Cache.Do("SETEX", sessionToken, "3600", userID)
	if err != nil {
		return http.Cookie{}, err
	}

	cookie := http.Cookie{
		Name:   "session_token",
		Value:  sessionToken,
		MaxAge: 3600,
	}

	return cookie, nil
}

func deleteAllTokensForUser(userID int) {
	keyData, err := conn.Cache.Do("KEYS", "*")
	if err != nil {
		log.Printf("Error getting all keys from cache: %s\n", err)
	}

	if keyData != nil {
		for _, k := range keyData.([]interface{}) {
			key := string([]byte(k.([]uint8)))

			sessUserID, err := conn.Cache.Do("GET", key)
			if err != nil {
				log.Printf("Error getting session token %s from cache\n", key)
				continue
			}
			if sessUserID == nil {
				continue
			}

			if userID == int(sessUserID.([]uint8)[0]) {
				_, err = conn.Cache.Do("DEL", key)
				if err != nil {
					log.Printf("Error deleting session token %s from cache: %s\n", key, err)
					continue
				}
			}
		}
	}
}

/* func generateToken(userID int, userName string) string {
	token := jwt.New(jwt.SigningMethodHS256)
	signingKey := os.Getenv("SIGNING_KEY")
	claims := token.Claims.(jwt.MapClaims)

	claims["userid"] = userID
	claims["username"] = userName

	tokenStr, _ := token.SignedString([]byte(signingKey))

	return tokenStr
} */
