// Package auth provides authentication functionalities using JWT tokens and cookies.
package auth

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v4"
	"go.uber.org/zap"
)

// tokenKey is the secret key used to sign the JWT tokens.
const tokenKey = "nextbug"

// Claims defines the structure of the JWT claims.
type Claims struct {
	jwt.RegisteredClaims
	UserID int `json:"user_id"`
}

// cryptoRandomID generates a cryptographically secure random ID within the given limit.
func cryptoRandomID(limit int) int {
	b := make([]byte, 8)
	_, err := rand.Read(b)
	if err != nil {
		return -1
	}
	return int(binary.BigEndian.Uint64(b) % uint64(limit))
}

// buildJWTString generates a new JWT token string with a random user ID.
func buildJWTString() (string, error) {
	id := cryptoRandomID(1000)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			// Token expiration time can be set here if needed.
			// ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenExp)),
			ExpiresAt: nil, // Token does not expire.
		},
		UserID: id,
	})
	tokenString, err := token.SignedString([]byte(tokenKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// getUserID extracts the user ID from the JWT token string.
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

// CheckCookie checks for a UserID cookie, creates one if it doesn't exist, and returns the user ID.
func CheckCookie(w http.ResponseWriter, r *http.Request, log *zap.Logger) (int, error) {
	var id int
	uuid, err := r.Cookie("UserID")
	if err != nil {
		if r.URL.Path == "/api/user/urls" {
			// Return a custom error with an informative message if the UserID cookie is missing.
			return 0, fmt.Errorf("cookie UserID is missing")
		}

		jwt, err := buildJWTString()
		if err != nil {
			log.Error("error making cookie", zap.Error(err))
			return 0, err
		}

		cookie := http.Cookie{
			Name:  "UserID",
			Value: jwt,
			Path:  "/",
		}

		// Set the cookie and return an error if it fails.
		http.SetCookie(w, &cookie)
		id, err = getUserID(jwt, log)
		if err != nil {
			log.Error("error creating cookie", zap.Error(err))
		}
		log.Info("generated UserID and token", zap.Int("UserID", id), zap.String("token key", jwt))
		return id, nil
	}

	id, err = getUserID(uuid.Value, log)
	if err != nil {
		log.Error("error parsing cookie", zap.Error(err))
	}
	log.Info("user ID retrieved from cookie", zap.Int("UserID", id))
	return id, nil
}
