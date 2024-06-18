package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hochitai/jpl/api/middleware"
	"github.com/hochitai/jpl/internal/model"
	"github.com/swaggo/swag/example/celler/httputil"
	"gorm.io/gorm"
)

type Word struct {}

func NewWordService() *Word {
	return &Word{}
}

func (w *Word) LoadAPIRouters(g *gin.RouterGroup, db *gorm.DB) {
	// Word
	g.GET("/words", w.GetWords(db))
	g.POST("/words", middleware.VerifyToken(), w.AddWord(db))
	g.PUT("/words/:id", middleware.VerifyToken(), w.UpdateWord(db))
	g.DELETE("/words/:id", middleware.VerifyToken(), w.DeleteWord(db))
	g.GET("/words/favorites", middleware.VerifyToken(), w.GetFavoriteWords(db))
}

// GetWords godoc
// @Summary      Get words
// @Description  Get words are pulic word
// @Tags         words
// @Accept       json
// @Produce      json
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success      200  {object}  []model.Word
// @Failure      500  {object}  httputil.HTTPError
// @Router       /v1/words [get]
func (w *Word) GetWords(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var wordModel model.Word
		words, err := wordModel.GetVocabularies(db)
		if err != nil {
			c.Error(err)
			// c.String(http.StatusInternalServerError, "Can not get vocabularies!")
			httputil.NewError(c, http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusOK, words)
	}
}

// GetFavoriteWords godoc
// @Summary      Get favorite words
// @Description  Get favorite words of current user
// @Tags         words
// @Accept       json
// @Produce      json
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success      200  {object}  []model.Word
// @Failure      500  {object}  httputil.HTTPError
// @Router       /v1/words/favorites [get]
func (w *Word) GetFavoriteWords(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userInfo := c.MustGet("userInfo").(model.User)
		favoriteModel := model.Favorite{UserId: userInfo.Id}
		words, err := favoriteModel.GetFavorites(db)
		if err != nil {
			c.Error(err)
			// c.String(http.StatusInternalServerError, "Can not get vocabularies!")
			httputil.NewError(c, http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusOK, words)
	}
}

// AddWord godoc
// @Summary      Add word
// @Description  Add word, if user is admin, word is public, else word is private
// @Tags         words
// @Accept       json
// @Produce      json
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param word body model.WordRequest true "object word"  Format(word)
// @Success      201  {object}  model.Word
// @Failure      400  {object}  httputil.HTTPError
// @Failure      500  {object}  httputil.HTTPError
// @Router       /v1/words [post]
func (w *Word) AddWord(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		decoder := json.NewDecoder(c.Request.Body)
		var word model.Word
		err := decoder.Decode(&word)
		if err != nil {
			c.Error(err)
			httputil.NewError(c ,http.StatusBadRequest, err)
			return
		}

		userInfo := c.MustGet("userInfo").(model.User)
		if userInfo.Permission == "admin" {
			word.Level = "public"
		} else {
			word.Level = "private"
		}

		err = word.CreateWord(db)
		if err != nil {
			c.Error(err)
			httputil.NewError(c ,http.StatusInternalServerError, fmt.Errorf("failure! Can not create new word"))
			return
		}

		if userInfo.Permission == "user" {
			favor := model.Favorite{WordId: word.Id, UserId: userInfo.Id}
			err = favor.CreateFavorite(db)
			if err != nil {
				c.Error(err)
				httputil.NewError(c ,http.StatusInternalServerError, fmt.Errorf("failure! Can not create new word"))
				return
			}
		}

		c.JSON(http.StatusCreated, word)
	}
}

// UpdateWord godoc
// @Summary      Update word
// @Description  Update word
// @Tags         words
// @Accept       json
// @Produce      json
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param 		 id 	path 	int		true 	"Word ID"
// @Param word body model.WordRequest true "object word"  Format(word)
// @Success      200  {object}  model.Word
// @Failure      400  {object}  httputil.HTTPError
// @Failure      500  {object}  httputil.HTTPError
// @Router       /v1/words/{id} [put]
func (w *Word) UpdateWord(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.Error(err)		
			httputil.NewError(c ,http.StatusBadRequest, fmt.Errorf("bad request body"))
			return
		}

		decoder := json.NewDecoder(c.Request.Body)
		var word model.Word
		err = decoder.Decode(&word)
		if err != nil {
			c.Error(err)	
			httputil.NewError(c ,http.StatusBadRequest, fmt.Errorf("bad request body"))
			return
		}
		word.SetId(id)

		userInfo := c.MustGet("userInfo").(model.User)
		favoriteModel := model.Favorite{UserId: userInfo.Id}
		words, err := favoriteModel.GetFavorites(db)
		if err != nil {
			c.Error(err)
			httputil.NewError(c, http.StatusInternalServerError, err)
			return
		}

		if userInfo.Permission == "user" {
			flag := false
			for _, u := range words {
				if u.Id == word.Id {
					flag = true
					break
				}
			}
			if !flag {
				httputil.NewError(c ,http.StatusInternalServerError,
					 fmt.Errorf("failure! You don't have permission to update this word"))
				return
			}
		}

		err = word.UpdateWord(db)
		if err != nil {
			c.Error(err)
			httputil.NewError(c ,http.StatusInternalServerError, fmt.Errorf("failure! Can not update word"))
			return
		}
		c.JSON(http.StatusOK, word)
	}
}

// DeleteWord godoc
// @Summary      Delete word
// @Description  Delete word and favorite word
// @Tags         words
// @Accept       json
// @Produce      json
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param 		 id 	path 	int		true 	"Word ID"
// @Success      200  string  string
// @Failure      400  {object}  httputil.HTTPError
// @Failure      500  {object}  httputil.HTTPError
// @Router       /v1/words/{id} [delete]
func (w *Word) DeleteWord(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.Error(err)	
			httputil.NewError(c ,http.StatusBadRequest, fmt.Errorf("bad request body"))
			return
		}
		var word model.Word
		word.SetId(id)

		err = word.DeleteWord(db)
		if err != nil {
			c.Error(err)
			httputil.NewError(c ,http.StatusInternalServerError, fmt.Errorf("failure! Can not delete word"))
			return
		}

		userInfo := c.MustGet("userInfo").(model.User)
		if userInfo.Permission == "user" {
			favor := model.Favorite{WordId: word.Id, UserId: userInfo.Id}
			err = favor.DeleteFavorite(db)
			if err != nil {
				c.Error(err)
				httputil.NewError(c ,http.StatusInternalServerError, fmt.Errorf("failure! Can not delete favorite word"))
				return
			}
		}
		c.String(http.StatusOK, "Delete successfully favorite word!")
	}
}