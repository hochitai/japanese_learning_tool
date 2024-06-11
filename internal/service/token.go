package service

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	firebase "firebase.google.com/go"
	"github.com/golang-jwt/jwt/v5"
	"github.com/hochitai/jpl/internal/model"
	"google.golang.org/api/option"
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

func CheckToken(tokenStr string) (model.User, error) {
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
	user, err := CheckToken(tokenStr);
	if err != nil {
		return "", err
	}
	token, err := CreateToken(user, time.Now().Add(time.Hour).Unix())
	if err != nil {
		return "", err
	}
	return token, nil
}

func CheckFirebaseToken(tokenStr string) (model.User, error) {
		var user model.User
		// Use a service account
		opt := option.WithCredentialsFile("serviceAccountKey.json")
		app, err := firebase.NewApp(context.Background(), nil, opt)
		if err != nil {
			return user, fmt.Errorf("token incorrect")
		}

		client, err := app.Auth(context.Background())
		if err != nil {
			return user, fmt.Errorf("token incorrect")
		}

		token, err := client.VerifyIDToken(context.Background(), tokenStr)
		if err != nil {
			return user, fmt.Errorf("token incorrect")
		}

		id, _ := strconv.Atoi(token.UID)
		user.Id = id
		user.Permission = "user"

		return user, nil
}