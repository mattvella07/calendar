package middleware

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	jwt "github.com/dgrijalva/jwt-go"
)

// Method contains the allowed http request methods for an endpoint
type Method struct {
	Allowed []string
}

// MethodChecker validates that the http request method is allowed for the endpoint
func (m Method) MethodChecker(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		for _, a := range m.Allowed {
			if a == r.Method {
				next.ServeHTTP(rw, r)
				return
			}
		}

		// HTTP request method is not allowed for the endpoint
		rw.WriteHeader(http.StatusMethodNotAllowed)
		rw.Write([]byte(fmt.Sprintf("HTTP request method must be one of the following %s for %s", m.Allowed, r.URL.String())))
	})
}

// ValidateJWT validates the JSON Web Token
func ValidateJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		token, _ := jwt.Parse(r.Header.Get("jwt"), func(token *jwt.Token) (interface{}, error) {
			signingKey := os.Getenv("SIGNING_KEY")

			return []byte(signingKey), nil
		})

		if token == nil || !token.Valid {
			log.Println("Invalid authorization token")
			rw.WriteHeader(http.StatusUnauthorized)
			rw.Write([]byte("Not Authorized"))
			return
		}

		// Set headers so endpoints can access JWT claims
		r.Header.Set("userid", strconv.Itoa(int(token.Claims.(jwt.MapClaims)["userid"].(float64))))
		r.Header.Set("username", token.Claims.(jwt.MapClaims)["username"].(string))

		log.Println("Valid token")
		next.ServeHTTP(rw, r)
	})
}
