package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
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

func AddAdmin(db *gorm.DB) gin.HandlerFunc {
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

		err = user.CreateAdmin(db)
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
		accessToken, err := service.CreateToken(user, time.Now().Add(time.Hour).Unix())
		if err != nil {
			c.String(http.StatusInternalServerError, "Error creating token")
		}
		user.AccessToken = accessToken

		// create refresh token
		refreshToken, err := service.CreateToken(user, time.Now().Add(24 * time.Hour).Unix())
		if err != nil {
			c.String(http.StatusInternalServerError, "Error creating token")
		}
		user.RefreshToken = refreshToken

		c.JSON(http.StatusOK, user)
	}
}

func UpdateUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.Error(err)		
			c.String(http.StatusBadRequest, "Bad request body")
			return
		}

		decoder := json.NewDecoder(c.Request.Body)
		var userLogin model.User
		err = decoder.Decode(&userLogin)
		if err != nil {
			c.Error(err)
			c.String(http.StatusBadRequest, "Bad request body")
			return
		}
		userLogin.Id = id

		userInfo := c.MustGet("userInfo").(model.User)
		if userInfo.Id == userLogin.Id  || userInfo.Permission == "admin" {
			err = userLogin.UpdateUser(db)
			if err != nil {
				c.Error(err)
				c.String(http.StatusInternalServerError, "Failure! Can not update user!")
				return
			}
			c.JSON(http.StatusOK, userLogin)
			return
		}
		c.String(http.StatusBadRequest, "You don't have permission to access this!")
	}
}

func DeleteUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.Error(err)		
			c.String(http.StatusBadRequest, "Bad request body")
			return
		}
		var user model.User
		user.Id = id

		err = user.DeleteUser(db)
		if err != nil {
			c.Error(err)
			c.String(http.StatusBadRequest, "Failure! Can not delete user!")
			return
		}
		c.JSON(http.StatusOK, user)

	}
}

func GetUsers(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var userModel model.User
		users, err := userModel.GetUsers(db)
		if err != nil {
			c.Error(err)
			c.String(http.StatusInternalServerError, "Can not get users!")
			return
		}
		c.JSON(http.StatusOK, users)
	}
}