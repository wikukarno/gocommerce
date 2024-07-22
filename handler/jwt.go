package handler

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("secret")

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func GenerateJWT(username string) (string, error) {
	expirationTime := time.Now().Add(1 * time.Minute)
	claims := &Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString , err := token.SignedString(jwtKey)
	
	if err != nil {
		return "", err
	}

	return tokenString, nil
	
}

func VerifyJWT(tokenString string) (*Claims, bool) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		return nil, false
	}

	return claims, token.Valid
}