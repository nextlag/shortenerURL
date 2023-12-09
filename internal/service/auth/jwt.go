package auth

import (
	"crypto/rand"
	"encoding/binary"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"go.uber.org/zap"
)

const SecretKey = "nextbug"
const TokenExp = time.Hour * 3

type Claims struct {
	jwt.RegisteredClaims
	UserID int `json:"user_id"`
}

func generateRandomID(limit int) int {
	b := make([]byte, 8)
	_, err := rand.Read(b)
	if err != nil {
		return -1
	}
	return int(binary.BigEndian.Uint64(b) % uint64(limit))
}

func buildJWTString() (string, error) {
	id := generateRandomID(1000)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TOKEN_EXP)),
		},
		UserID: id,
	})
	tokenString, err := token.SignedString([]byte(SECRET_KEY))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func getUserID(tokenString string, log *zap.Logger) int {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Error("unexpected signing method")
			return nil, nil
		}
		return []byte(SECRET_KEY), nil
	})
	if err != nil || !token.Valid {
		log.Error("Token is not valid", zap.Error(err))
		return -1
	}
	return claims.UserID
}

func CheckCookieForID(w http.ResponseWriter, r *http.Request, log *zap.Logger) int {
	var id int
	userIDCookie, err := r.Cookie("UserID")
	if err != nil {
		if r.URL.Path == "/api/user/urls" {
			return -1
		}
		jwt, err := buildJWTString()
		if err != nil {
			log.Error("Error making cookie", zap.Error(err))
			return -1
		}
		cookie := http.Cookie{Name: "UserID", Value: jwt}
		http.SetCookie(w, &cookie)
		id = getUserID(jwt, log)
		return id
	}
	id = getUserID(userIDCookie.Value, log)
	return id
}
