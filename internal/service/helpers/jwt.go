package helpers

import (
	"errors"
	"github.com/golang-jwt/jwt"
	"gitlab.com/distributed_lab/Auth/internal/data"
)

func GenerateAccessToken(user data.User, expires int64, secret string, permissions string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = expires
	claims["owner_id"] = user.Id
	claims["email"] = user.Email
	claims["module.permission"] = permissions

	return token.SignedString([]byte(secret))
}

func GenerateRefreshToken(user data.User, expires int64, secret string, permissions string) (string, error, jwt.MapClaims) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = expires
	claims["owner_id"] = user.Id
	claims["email"] = user.Email
	claims["module.permission"] = permissions

	signedToken, err := token.SignedString([]byte(secret))

	return signedToken, err, claims
}

func CheckRefreshToken(tokenStr string, ownerId int64, secret string) error {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}

		return []byte(secret), nil
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

func ParseJwtToken(tokenStr string, secret string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}

		return []byte(secret), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, err
}
