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
			c.String(http.StatusInternalServerError, "Can not get vocabulary!")
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

		err = word.CreateWord(db)
		if err != nil {
			c.Error(err)
			c.String(http.StatusInternalServerError, "Failure! Can not create new word!")
			return
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
			c.String(http.StatusInternalServerError, "Failure! Can not update new word!")
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
		decoder := json.NewDecoder(c.Request.Body)
		var word model.Word
		err = decoder.Decode(&word)
		if err != nil {
			c.Error(err)	
			c.String(http.StatusBadRequest, "Bad request body")
			return
		}
		word.SetId(id)

		err = word.DeleteWord(db)
		if err != nil {
			c.Error(err)
			c.String(http.StatusBadRequest, "Failure! Can not delete new word!")
			return
		}
		c.JSON(http.StatusOK, word)
	}
}