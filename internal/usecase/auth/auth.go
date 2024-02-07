package auth

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v4"
	"go.uber.org/zap"
)

// const tokenExp = time.Hour * 3
const tokenKey = "nextbug"

type Claims struct {
	jwt.RegisteredClaims
	UserID int `json:"user_id"`
}

func cryptoRandomID(limit int) int {
	b := make([]byte, 8)
	_, err := rand.Read(b)
	if err != nil {
		return -1
	}
	return int(binary.BigEndian.Uint64(b) % uint64(limit))
}

func buildJWTString() (string, error) {
	id := cryptoRandomID(1000)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			// хранение токена по времени
			// ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenExp)),

			// хранение токена бессрочно
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

func CheckCookie(w http.ResponseWriter, r *http.Request, log *zap.Logger) (int, error) {
	var id int
	uuid, err := r.Cookie("UserID")
	if err != nil {
		if r.URL.Path == "/api/user/urls" {
			// Возвращаем пользовательскую ошибку с информативным сообщением
			return 0, fmt.Errorf("файл cookie UserID отсутствует")
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

		// Возвращаем ошибку при установке куки, если она не удастся
		http.SetCookie(w, &cookie)
		id, err = getUserID(jwt, log)
		if err != nil {
			log.Error("ошибка создания файла cookie", zap.Error(err))
		}
		log.Info("generate UserID and token:", zap.Int("UserID", id), zap.String("token key", jwt))
		return id, nil
	}

	id, err = getUserID(uuid.Value, log)
	if err != nil {
		log.Error("ошибка создания файла cookie", zap.Error(err))
	}
	log.Info("ID пользователя после получения из cookie", zap.Int("UserID", id))
	return id, nil
}
