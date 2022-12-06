package helpers

import (
	"errors"
	"github.com/golang-jwt/jwt"
	"github.com/mhrynenko/jwt_service/internal/data"
	"time"
)

const tmpSecret = "my secret key"

func GenerateAccessToken(user data.User) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(10 * time.Minute).Unix()
	claims["owner_id"] = user.Id
	claims["email"] = user.Email

	return token.SignedString([]byte(tmpSecret))
}

func GenerateRefreshToken(user data.User) (string, error, jwt.MapClaims) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(24 * time.Hour).Unix()
	claims["owner_id"] = user.Id
	claims["email"] = user.Email

	signedToken, err := token.SignedString([]byte(tmpSecret))

	return signedToken, err, claims
}

func CheckRefreshToken(tokenStr string, ownerId int64) error {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}

		return []byte(tmpSecret), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if int64(claims["owner_id"].(float64)) != ownerId {
			return errors.New("invalid token")
		}
	}

	if !token.Valid {
		return errors.New("invalid token")
	}

	return err
}
