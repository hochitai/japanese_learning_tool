package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hochitai/jpl/internal/model"
	"gorm.io/gorm"
)


func GetWords(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var wordModel model.Word
		words, err := wordModel.GetVocabularies(db)
		if err != nil {
			c.Error(err)
			c.String(http.StatusInternalServerError, "Can not get vocabularies!")
			return
		}
		c.JSON(http.StatusOK, words)
	}

}

func AddWord(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		decoder := json.NewDecoder(c.Request.Body)
		var word model.Word
		err := decoder.Decode(&word)
		if err != nil {
			c.Error(err)
			c.String(http.StatusBadRequest, "Bad request body")
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
			c.String(http.StatusInternalServerError, "Failure! Can not create new word!")
			return
		}

		if userInfo.Permission == "user" {
			favor := model.Favorite{WordId: word.Id, UserId: userInfo.Id}
			err = favor.CreateFavorite(db)
			if err != nil {
				c.Error(err)
				c.String(http.StatusInternalServerError, "Failure! Can not create new word!")
				return
			}
		}

		c.JSON(http.StatusCreated, word)
	}
}

func UpdateWord(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.Error(err)		
			c.String(http.StatusBadRequest, "Bad request body")
			return
		}

		decoder := json.NewDecoder(c.Request.Body)
		var word model.Word
		err = decoder.Decode(&word)
		if err != nil {
			c.Error(err)	
			c.String(http.StatusBadRequest, "Bad request body")
			return
		}
		word.SetId(id)

		err = word.UpdateWord(db)
		if err != nil {
			c.Error(err)
			c.String(http.StatusInternalServerError, "Failure! Can not update word!")
			return
		}
		c.JSON(http.StatusOK, word)
	}
}

func DeleteWord(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.Error(err)	
			c.String(http.StatusBadRequest, "Bad request body")
			return
		}
		var word model.Word
		word.SetId(id)

		err = word.DeleteWord(db)
		if err != nil {
			c.Error(err)
			c.String(http.StatusBadRequest, "Failure! Can not delete word!")
			return
		}

		userInfo := c.MustGet("userInfo").(model.User)
		if userInfo.Permission == "user" {
			favor := model.Favorite{WordId: word.Id, UserId: userInfo.Id}
			err = favor.DeleteFavorite(db)
			if err != nil {
				c.Error(err)
				c.String(http.StatusInternalServerError, "Failure! Can not delete favorite word!")
				return
			}
		}

		c.String(http.StatusOK, "Delete successfully favorite word!")
	}
}