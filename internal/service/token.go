package service

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/hochitai/jpl/internal/model"
)

func CreateToken(user model.User, duration int64) (string, error) {
	secretKey := []byte(os.Getenv("SECRET_KEY"))

	// Create a new token object, specifying signing method and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.Id,                         // Subject (user identifier)
		"per": user.Permission,					// Permission
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

func ValidateToken(tokenStr string) (model.User, error) {
	secretKey := []byte(os.Getenv("SECRET_KEY"))
	var user model.User
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// secretKey is a []byte containing your secret
		return secretKey, nil
	})
	if err != nil {
		return user, fmt.Errorf("token incorrect")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		user.Id = int(claims["sub"].(float64))
		user.Permission = claims["per"].(string)
	} 
	return user, nil
}

func RefreshToken(tokenStr string) (string, error) {
	user, err := ValidateToken(tokenStr);
	if err != nil {
		return "", err
	}
	token, err := CreateToken(user, time.Now().Add(time.Hour).Unix())
	if err != nil {
		return "", err
	}
	return token, nil
}