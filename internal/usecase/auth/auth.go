// Package auth provides authentication utilities, including JWT generation and validation,
// and functions to manage user ID cookies in HTTP requests.
package auth

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v4"
	"go.uber.org/zap"
)

// tokenKey is the key used for signing JWT tokens.
const tokenKey = "nextbug"

// Claims represents the structure of JWT claims.
type Claims struct {
	jwt.RegisteredClaims
	UserID int `json:"user_id"`
}

// cryptoRandomID generates a cryptographically secure random integer ID up to the specified limit.
func cryptoRandomID(limit int) int {
	b := make([]byte, 8)
	_, err := rand.Read(b)
	if err != nil {
		return -1
	}
	return int(binary.BigEndian.Uint64(b) % uint64(limit))
}

// buildJWTString creates a new JWT string with a random user ID.
func buildJWTString() (string, error) {
	id := cryptoRandomID(1000)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			// Token storage duration
			// ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenExp)),

			// Token stored indefinitely
			ExpiresAt: nil,
		},
		UserID: id,
	})
	tokenString, err := token.SignedString([]byte(tokenKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// getUserID parses the JWT token string and retrieves the user ID from its claims.
func getUserID(tokenString string, log *zap.Logger) (int, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Error("unexpected signing method")
			return nil, nil
		}
		return []byte(tokenKey), nil
	})
	if err != nil || !token.Valid {
		log.Error("Token is not valid", zap.Error(err))
		return 0, err
	}
	return claims.UserID, nil
}

// CheckCookie checks for a user ID cookie in the request. If not found, it generates a new one.
// It returns the user ID and any error encountered.
func CheckCookie(w http.ResponseWriter, r *http.Request, log *zap.Logger) (int, error) {
	var id int
	uuid, err := r.Cookie("UserID")
	if err != nil {
		if r.URL.Path == "/api/user/urls" {
			// Return a user-friendly error if the cookie is missing.
			return 0, fmt.Errorf("UserID cookie is missing")
		}

		var jwt string
		jwt, err = buildJWTString()
		if err != nil {
			log.Error("error making cookie", zap.Error(err))
			return 0, err
		}

		cookie := http.Cookie{
			Name:  "UserID",
			Value: jwt,
			Path:  "/",
		}

		// Set the cookie, returning an error if unsuccessful.
		http.SetCookie(w, &cookie)
		id, err = getUserID(jwt, log)
		if err != nil {
			log.Error("error creating UserID cookie", zap.Error(err))
		}
		log.Info("Generated UserID and token", zap.Int("UserID", id), zap.String("token key", jwt))
		return id, nil
	}

	id, err = getUserID(uuid.Value, log)
	if err != nil {
		log.Error("error retrieving UserID from cookie", zap.Error(err))
	}
	log.Info("UserID retrieved from cookie", zap.Int("UserID", id))
	return id, nil
}
