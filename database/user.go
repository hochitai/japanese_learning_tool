package database

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dchest/uniuri"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	_ "github.com/joho/godotenv/autoload"
	"gorm.io/gorm"
)

type User struct {
	Id 			int `json:"id"`
	Username 	string `json:"username"`
	Password 	string `json:"password,omitempty"`
	Name		string `json:"name"`
	Salt 		string `json:"-"`
	Token 		string `json:"token" gorm:"-"` // token with duration 1h
}

func AddUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		decoder := json.NewDecoder(c.Request.Body)
		var user User
		err := decoder.Decode(&user)
		if err != nil {
			c.Error(err)
			c.String(http.StatusBadRequest, "Bad request body")
			return
		}

		// Create random string
		salt := uniuri.New()
		hashedPassword := hashPassword(user.Password, salt)

		// Update password and salt 
		user.Password = hashedPassword
		user.Salt = salt

		result := db.Create(&user)
		if result.Error != nil {
			c.Error(result.Error)
			c.String(http.StatusBadRequest, "Failue to create new user!")
			return
		}

		c.JSON(http.StatusOK, user)
		
	}
}

func hashPassword(password string, salt string) string {
	key := os.Getenv("KEY")

	// Hash password with format key + password + salt
	h := md5.New()
	h.Write([]byte(key+password+salt))
	bytesPassword := h.Sum(nil)
	hashedPassword := hex.EncodeToString(bytesPassword)
	return hashedPassword
}

func Login(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		decoder := json.NewDecoder(c.Request.Body)
		var userLogin User
		err := decoder.Decode(&userLogin)
		if err != nil {
			c.Error(err)
			c.String(http.StatusBadRequest, "Bad request body")
			return
		}

		var user User
		result := db.Where("username = ?",userLogin.Username).First(&user)
		if result.Error != nil {
			c.Error(result.Error)
			c.String(http.StatusBadRequest, "Bad request body")
			return
		}
		
		hashedPassword := hashPassword(userLogin.Password, user.Salt)

		if strings.Compare(hashedPassword, user.Password) != 0 {
			c.String(http.StatusBadRequest, "Username or password is incorrect!")
			return
		}
		// Hidden field password
		user.Password = ""
		
		token, err := createToken(user.Name)
		if err != nil {
			c.String(http.StatusInternalServerError, "Error creating token")
		}

		user.Token = token

		c.JSON(http.StatusOK, user)
	}
}

func createToken(username string) (string, error) {
	secretKey := []byte(os.Getenv("SECRET_KEY"))

	// Create a new token object, specifying signing method and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": username,                    		// Subject (user identifier)
		"iss": "login",                  			// Issuer        
		"exp": time.Now().Add(30 * time.Second).Unix(), 	// Expiration time
		"iat": time.Now().Unix(),  					// Issued at
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", nil
	}
	return tokenString, nil
}

func validateToken(tokenStr string) (bool, error) {
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
		return false, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if t, ok := claims["exp"].(float64); ok {
			if float64(time.Now().Unix()) - t < 0 {
				return true, nil
			}
			return false, fmt.Errorf("token was expired")
		}
	} else {
		return false, fmt.Errorf("incorrect Token provided")
	}
	return true, nil
}

func UpdateUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		
	}
}

func DeleteUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		
	}
}

func GetUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		
	}
}