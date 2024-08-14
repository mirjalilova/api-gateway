package tokens

import (
	"errors"
	"fmt"
	"golang.org/x/exp/slog"

	"github.com/golang-jwt/jwt"
)

func VerifyToken(tokenString string) error {
	secretKey := []byte("secret")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		slog.Error("Error while parsing token: ", err)
		return err
	}

	if !token.Valid {
		slog.Error("Invalid token")
		return fmt.Errorf("invalid token")
	}

	return nil
}

func ExtractClaims(tokenStr, key string) (jwt.MapClaims, error) {
	var (
		token *jwt.Token
		err   error
	)

	keyFunc := func(token *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	}
	token, err = jwt.Parse(tokenStr, keyFunc)
	if err != nil {
		slog.Error("Error while parsing token: ", err)
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !(ok && token.Valid) {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
