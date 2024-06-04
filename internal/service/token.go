package service

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func CreateToken(username string, duration int64) (string, error) {
	secretKey := []byte(os.Getenv("SECRET_KEY"))

	// Create a new token object, specifying signing method and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": username,                         // Subject (user identifier)
		"iss": "login",                          // Issuer
		"exp": duration, // Expiration time
		"iat": time.Now().Unix(),                // Issued at
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", nil
	}
	return tokenString, nil
}

func ValidateToken(tokenStr string) (string, error) {
	secretKey := []byte(os.Getenv("SECRET_KEY"))
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// secretKey is a []byte containing your secret
		return secretKey, nil
	})
	if err != nil {
		return "", fmt.Errorf("token incorrect")
	}

	var username string
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		username = claims["sub"].(string)
	} 
	return username, nil
}

func RefreshToken(tokenStr string) (string, error) {
	username, err := ValidateToken(tokenStr);
	if err != nil {
		return "", err
	}
	token, err := CreateToken(username, time.Now().Add(time.Hour).Unix())
	if err != nil {
		return "", err
	}
	return token, nil
}