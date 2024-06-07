package model

import (
	"math/rand"
	"strings"
	"time"

	"github.com/hochitai/jpl/internal/database"
	"gorm.io/gorm"
)

type Word struct {
	Id            int    `gorm:"primaryKey" json:"id"`
	Characters    string `json:"characters"`
	Pronunciation string `json:"pronunciation"`
	Meaning       string `json:"meaning"`
	Level 		  string `json:"-"`
}

func (w *Word) SetId(id int) {
	w.Id = id
}

func (word *Word) CreateWord(db *gorm.DB) error {
	result := db.Create(&word)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (word *Word) UpdateWord(db *gorm.DB) error {
	result := db.Omit("level").Save(&word)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (word *Word) DeleteWord(db *gorm.DB) error {
	result := db.Delete(&word)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (word Word) GetVocabularies(db *gorm.DB) ([]Word, error) {
	var words []Word
	result := db.Where("level = ?", "public").Find(&words)
	if result.Error != nil {
		return nil, result.Error
	}
	return words, nil
}

func GetHiraganaCharacters() []Word {
	return getCharacters(database.HiraganaTable)
}

func GetKatakanaCharacters() []Word {
	return getCharacters(database.KatakanaTable)
}

func GetRandomCharacter(words []Word) Word {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	max := len(words)
	index := r.Intn(max)
	return words[index]
}

func getCharacters(alphabet string) []Word {
	var words []Word
	rows := strings.Split(alphabet, "\n")
	for _, row := range rows[1:] {
		phrases := strings.Split(row, "\t")[1:]
		for _, phrase := range phrases {
			word := strings.Split(phrase, " ")
			if len(word) > 1 {
				words = append(words, Word{
					Characters:    word[0],
					Pronunciation: word[1],
				})
			}
		}
	}
	return words
}