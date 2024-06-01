package database

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"github.com/dchest/uniuri"
	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v10"
	"github.com/joho/godotenv"
)

type User struct {
	Id 			int `json:"id"`
	Username 	string `json:"username"`
	Password 	string `json:"password,omitempty"`
	Name		string `json:"name"`
	Salt 		string `json:"-"`
}

// func (u *User) SetPassword(password string) {
// 	u.Password = password
// }

// func (u User) GetPassword() string {
// 	return u.Password
// }

// func (u *User) SetSalt(salt string) {
// 	u.salt = salt
// }

func AddUser(db *pg.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := godotenv.Load()
		if err != nil {
			c.Error(err)
			c.JSON(http.StatusBadRequest, struct{}{})
			return
		}
		key := os.Getenv("KEY")

		decoder := json.NewDecoder(c.Request.Body)
		var user User
		err = decoder.Decode(&user)
		if err != nil {
			c.Error(err)
			c.JSON(http.StatusBadRequest, struct{}{})
			return
		}

		// Create random string
		salt := uniuri.New()

		// Hash password with format key + password + salt
		h := md5.New()
		h.Write([]byte(key+user.Password+salt))
		bytesPassword := h.Sum(nil)
		hashedPassword := hex.EncodeToString(bytesPassword)

		// Update password and salt 
		user.Password = hashedPassword
		user.Salt = salt

		_, err = db.Model(&user).Returning("*").Insert()
		if err != nil {
			c.Error(err)
			c.JSON(http.StatusBadRequest, struct{}{})
			return
		}

		c.JSON(http.StatusOK, user)
		
	}
}

func UpdateUser(db *pg.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		
	}
}

func DeleteUser(db *pg.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		
	}
}

func GetUser(db *pg.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		
	}
}

func Login(db *pg.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := godotenv.Load()
		if err != nil {
			c.Error(err)
			c.JSON(http.StatusBadRequest, struct{}{})
			return
		}
		key := os.Getenv("KEY")

		decoder := json.NewDecoder(c.Request.Body)
		var userLogin User
		err = decoder.Decode(&userLogin)
		if err != nil {
			c.Error(err)
			c.JSON(http.StatusBadRequest, struct{}{})
			return
		}

		var user User
		err = db.Model(&user).Where("username = ?", userLogin.Username).Select()
		
		// Hash password with format key + password + salt
		h := md5.New()
		h.Write([]byte(key+userLogin.Password+user.Salt))
		bytesPassword := h.Sum(nil)
		hashedPassword := hex.EncodeToString(bytesPassword)

		if err != nil {
			c.Error(err)
			c.JSON(http.StatusBadRequest, struct{}{})
			return
		}

		if strings.Compare(hashedPassword, user.Password) != 0 {
			c.JSON(http.StatusBadRequest, struct{}{})
			return
		}

		user.Password = ""
		c.JSON(http.StatusOK, user)
	}
}