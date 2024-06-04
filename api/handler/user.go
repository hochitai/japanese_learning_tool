package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dchest/uniuri"
	"github.com/gin-gonic/gin"
	"github.com/hochitai/jpl/internal/model"
	"github.com/hochitai/jpl/internal/service"
	_ "github.com/joho/godotenv/autoload"
	"gorm.io/gorm"
)

func AddUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		decoder := json.NewDecoder(c.Request.Body)
		var user model.User
		err := decoder.Decode(&user)
		if err != nil {
			c.Error(err)
			c.String(http.StatusBadRequest, "Bad request body")
			return
		}

		// Create random string
		salt := uniuri.New()
		hashedPassword := service.HashPassword(user.Password, salt)

		// Update password and salt 
		user.Password = hashedPassword
		user.Salt = salt

		err = user.CreateUser(db)
		if err != nil {
			c.Error(err)
			c.String(http.StatusBadRequest, "Failue to create new user!")
			return
		}

		c.JSON(http.StatusOK, user)
		
	}
}

func Login(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		decoder := json.NewDecoder(c.Request.Body)
		var userLogin model.User
		err := decoder.Decode(&userLogin)
		if err != nil {
			c.Error(err)
			c.String(http.StatusBadRequest, "Bad request body")
			return
		}

		user, err := userLogin.GetUserByUsername(db)
		if err != nil {
			c.Error(err)
			c.String(http.StatusBadRequest, "Bad request body")
			return
		}
		
		hashedPassword := service.HashPassword(userLogin.Password, user.Salt)

		if strings.Compare(hashedPassword, user.Password) != 0 {
			c.String(http.StatusBadRequest, "Username or password is incorrect!")
			return
		}
		// Hidden field password
		user.Password = ""
		
		// create access token
		accessToken, err := service.CreateToken(user.Name, time.Now().Add(time.Hour).Unix())
		if err != nil {
			c.String(http.StatusInternalServerError, "Error creating token")
		}
		user.AccessToken = accessToken

		// create refresh token
		refreshToken, err := service.CreateToken(user.Name, time.Now().Add(24 * time.Hour).Unix())
		if err != nil {
			c.String(http.StatusInternalServerError, "Error creating token")
		}
		user.RefreshToken = refreshToken

		c.JSON(http.StatusOK, user)
	}
}

func UpdateUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("Update...")
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