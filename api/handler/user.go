package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/dchest/uniuri"
	"github.com/gin-gonic/gin"
	"github.com/hochitai/jpl/api/middleware"
	"github.com/hochitai/jpl/internal/model"
	"github.com/hochitai/jpl/internal/service"
	_ "github.com/joho/godotenv/autoload"
	"github.com/swaggo/swag/example/celler/httputil"
	"gorm.io/gorm"
)

type User struct {}

func NewUserService() *User {
	return &User{}
}

func (u *User) LoadAPIRouters(g *gin.RouterGroup, db *gorm.DB) {
	// User
	g.POST("/users/register", u.AddUser(db))
	g.POST("/users/admin", middleware.VerifyTokenAndPermission(), u.AddAdmin(db))
	g.POST("/users/login", u.Login(db))
	g.POST("/users/token", middleware.RefreshToken())
	g.PUT("/users/:id", middleware.VerifyToken(), u.UpdateUser(db))
	g.DELETE("/users/:id", middleware.VerifyTokenAndPermission(), u.DeleteUser(db))

	//Admin
	g.GET("/admin/users", middleware.VerifyTokenAndPermission(), u.GetUsers(db))
}

// AddUser godoc
// @Summary      Add user
// @Description  Add user (register)
// @Tags         users
// @Accept       json
// @Produce      json
// @Param user body model.UserInfoRequest true "object user"  Format(user)
// @Success      201  {object}  model.UserResponse
// @Failure      400  {object}  httputil.HTTPError
// @Router       /v1/users/register [post]
func (u *User) AddUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		decoder := json.NewDecoder(c.Request.Body)
		var user model.User
		err := decoder.Decode(&user)
		if err != nil {
			c.Error(err)
			httputil.NewError(c ,http.StatusBadRequest, fmt.Errorf("bad request body"))
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
			httputil.NewError(c ,http.StatusBadRequest, fmt.Errorf("bad request body"))
			return
		}

		c.JSON(http.StatusCreated, model.UserResponse{Id: user.Id, Name: user.Name})
		
	}
}

// AddAdmin godoc
// @Summary      Add admin
// @Description  Add admin, only admin can add admin
// @Tags         users
// @Accept       json
// @Produce      json
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param user body model.UserInfoRequest true "object user"  Format(user)
// @Success      201  {object}  model.UserResponse
// @Failure      400  {object}  httputil.HTTPError
// @Router       /v1/users/admin [post]
func (u *User) AddAdmin(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		decoder := json.NewDecoder(c.Request.Body)
		var user model.User
		err := decoder.Decode(&user)
		if err != nil {
			c.Error(err)
			httputil.NewError(c ,http.StatusBadRequest, fmt.Errorf("bad request body"))
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
			httputil.NewError(c ,http.StatusInternalServerError, fmt.Errorf("failue to create new user"))
			return
		}

		c.JSON(http.StatusCreated, model.UserResponse{Id: user.Id, Name: user.Name})
		
	}
}

// Login godoc
// @Summary      Login
// @Description  Login
// @Tags         users
// @Accept       json
// @Produce      json
// @Param user body model.UserLoginRequest true "object user" Format(user)
// @Success      200  {object}  model.UserLoginResponse
// @Failure      400  {object}  httputil.HTTPError
// @Failure      500  {object}  httputil.HTTPError
// @Router       /v1/users/login [post]
func (u *User) Login(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		decoder := json.NewDecoder(c.Request.Body)
		var userLogin model.User
		err := decoder.Decode(&userLogin)
		if err != nil {
			c.Error(err)
			httputil.NewError(c ,http.StatusBadRequest, fmt.Errorf("bad request body"))
			return
		}

		user, err := userLogin.GetUserByUsername(db)
		if err != nil {
			c.Error(err)
			httputil.NewError(c ,http.StatusBadRequest, fmt.Errorf("bad request body"))
			return
		}
		
		hashedPassword := service.HashPassword(userLogin.Password, user.Salt)

		if strings.Compare(hashedPassword, user.Password) != 0 {
			httputil.NewError(c ,http.StatusInternalServerError, fmt.Errorf("username or password is incorrect"))
			return
		}
		// Hidden field password
		user.Password = ""
		
		// create access token
		accessToken, err := service.CreateToken(user, time.Now().Add(time.Hour).Unix())
		if err != nil {
			httputil.NewError(c ,http.StatusInternalServerError, fmt.Errorf("error creating token"))
		}
		user.AccessToken = accessToken

		// create refresh token
		refreshToken, err := service.CreateToken(user, time.Now().Add(24 * time.Hour).Unix())
		if err != nil {
			httputil.NewError(c ,http.StatusInternalServerError, fmt.Errorf("error creating token"))
		}
		user.RefreshToken = refreshToken

		c.JSON(http.StatusOK, model.UserLoginResponse{
				Name: user.Name, 
				AccessToken: user.AccessToken, 
				RefreshToken: user.RefreshToken,
			})
	}
}

// UpdateUser godoc
// @Summary      Update user
// @Description  Update user
// @Tags         users
// @Accept       json
// @Produce      json
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param 		 id 	path 	int		true 	"User ID"
// @Param user body model.UserRequest true "object user" Format(user)
// @Success      200  {object}  model.UserResponse
// @Failure      400  {object}  httputil.HTTPError
// @Failure      500  {object}  httputil.HTTPError
// @Router       /v1/users/{id} [put]
func (u *User) UpdateUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.Error(err)		
			httputil.NewError(c ,http.StatusBadRequest, fmt.Errorf("bad request body"))
			return
		}

		decoder := json.NewDecoder(c.Request.Body)
		var userLogin model.User
		err = decoder.Decode(&userLogin)
		if err != nil {
			c.Error(err)
			httputil.NewError(c ,http.StatusBadRequest, fmt.Errorf("bad request body"))
			return
		}
		userLogin.Id = id

		userInfo := c.MustGet("userInfo").(model.User)
		if userInfo.Id == userLogin.Id  || userInfo.Permission == "admin" {
			err = userLogin.UpdateUser(db)
			if err != nil {
				c.Error(err)
				httputil.NewError(c ,http.StatusInternalServerError, fmt.Errorf("failure! Can not update user"))
				return
			}
			c.JSON(http.StatusOK, model.UserResponse{Id: userLogin.Id, Name: userLogin.Name})
			return
		}
		httputil.NewError(c ,http.StatusBadRequest, fmt.Errorf("you don't have permission to access this"))
	}
}

// DeleteUser godoc
// @Summary      Delete user
// @Description  Delete user
// @Tags         users
// @Accept       json
// @Produce      json
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param 		 id 	path 	int		true 	"User ID"
// @Success      200  string  string
// @Failure      400  {object}  httputil.HTTPError
// @Failure      500  {object}  httputil.HTTPError
// @Router       /v1/users/{id} [delete]
func (u *User) DeleteUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.Error(err)		
			httputil.NewError(c ,http.StatusBadRequest, fmt.Errorf("bad request body"))
			return
		}
		var user model.User
		user.Id = id

		err = user.DeleteUser(db)
		if err != nil {
			c.Error(err)
			httputil.NewError(c ,http.StatusInternalServerError, fmt.Errorf("failure! Can not delete user"))
			return
		}
		c.String(http.StatusOK, "Delete successfully favorite word!")

	}
}

// GetUsers godoc
// @Summary      Get users
// @Description  Get users if you are admin
// @Tags         users
// @Accept       json
// @Produce      json
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success      200  {object}  []model.UserResponse
// @Failure      500  {object}  httputil.HTTPError
// @Router       /v1/admin/users [get]
func (u *User) GetUsers(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var userModel model.User
		users, err := userModel.GetUsers(db)
		if err != nil {
			c.Error(err)
			c.String(http.StatusInternalServerError, "Can not get users!")
			return
		}
		userResponse := make([]model.UserResponse, 0)
		for _, user := range users {
			userResponse = append(userResponse, model.UserResponse{Id: user.Id, Name: user.Name})
		}
		c.JSON(http.StatusOK, userResponse)
	}
}

func (u *User) HomePage(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "auth.tmpl", gin.H{
			"title": "Main website",
		})
	}
}