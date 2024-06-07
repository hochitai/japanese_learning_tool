package middleware

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hochitai/jpl/internal/model"
	"github.com/hochitai/jpl/internal/service"
)

func CheckToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := strings.Split(c.Request.Header.Get("Authorization"), " ")[1]
		user, err :=  service.ValidateToken(token)
		if err != nil {
			c.String(http.StatusUnauthorized, err.Error())
			return
		}
		c.Set("userInfo", user)
		c.Next()
	}
}

func RefreshToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		// client send refresh-token in request body
		decoder := json.NewDecoder(c.Request.Body)
		var userLogin model.User
		err := decoder.Decode(&userLogin)
		if err != nil {
			c.Error(err)
			c.String(http.StatusBadRequest, "Bad request body")
			return
		}

		accessToken, err :=  service.RefreshToken(userLogin.RefreshToken)
		if err != nil {
			c.String(http.StatusUnauthorized, err.Error())
			return
		}

		userLogin.AccessToken = accessToken

		c.JSON(http.StatusOK, userLogin)
	}
}

func CheckTokenAndPermission() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := strings.Split(c.Request.Header.Get("Authorization"), " ")[1]
		user, err :=  service.ValidateToken(token)
		if err != nil || user.Permission != "admin" {
			c.String(http.StatusUnauthorized, "You don't have permission to access this!")
			return
		}
		c.Set("userInfo", user)
		c.Next()
	}
}