package helper

import (
	jwt "github.com/dgrijalva/jwt-go"
)

func GetToken(id int, email, role string) (string, error) {
	signingKey := []byte("rahasia")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    id,
		"email": email,
		"role":  role,
	})
	tokenString, err := token.SignedString(signingKey)
	return tokenString, err
}

func VerifyToken(tokenString string) (jwt.Claims, error) {
	signingKey := []byte("rahasia")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return signingKey, nil
	})

	if err != nil {
		return nil, err
	}
	return token.Claims, err
}
