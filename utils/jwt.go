package utils

import (
	"errors"
	"time"

	"example.com/config"
	"github.com/golang-jwt/jwt/v5"
)

var accessSecretKey string
var refreshSecretKey string

func init() {
	accessSecretKey = config.GetEnv("ACCESS_TOKEN_SECRET")
	if accessSecretKey == "" {
		panic("ACCESS_TOKEN_SECRET not set in environment")
	}
	refreshSecretKey = config.GetEnv("REFRESH_TOKEN_SECRET")
	if refreshSecretKey == "" {
		panic("REFRESH_TOKEN_SECRET not set in environment")
	}
}

func GenerateAccessToken(login string, userId int64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"login":  login,
		"userId": userId,
		"exp":    time.Now().Add(time.Minute * 15).Unix(),
	})

	return token.SignedString([]byte(accessSecretKey))
}

func VerifyAccessToken(token string) (int64, error) {
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)

		if !ok {
			return nil, errors.New("unexpected signing method")
		}

		return []byte(accessSecretKey), nil
	})

	if err != nil {
		return 0, errors.New("could not parse token")
	}

	tokenIsValid := parsedToken.Valid

	if !tokenIsValid {
		return 0, errors.New("invalid token")
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)

	if !ok {
		return 0, errors.New("invalid token claims")
	}

	// login := claims["login"].(string) // login,  ok := ...
	userId := int64(claims["userId"].(float64))

	return userId, nil

}

func GenerateRefreshToken(login string, userId int64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"login":  login,
		"userId": userId,
		"exp":    time.Now().Add(time.Hour * 48).Unix(),
	})

	return token.SignedString([]byte(refreshSecretKey))
}

func VerifyRefreshToken(token string) (int64, string, error) {
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)

		if !ok {
			return nil, errors.New("unexpected signing method")
		}

		return []byte(refreshSecretKey), nil
	})

	if err != nil {
		return 0, "", errors.New("could not parse token")
	}

	tokenIsValid := parsedToken.Valid

	if !tokenIsValid {
		return 0, "", errors.New("invalid token")
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)

	if !ok {
		return 0, "", errors.New("invalid token claims")
	}

	login := claims["login"].(string)
	userId := int64(claims["userId"].(float64))

	return userId, login, nil

}
