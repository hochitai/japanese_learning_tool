package model

import "gorm.io/gorm"

type Favorite struct {
	Id     int `json:"id"`
	WordId int `json:"word_id"`
	UserId int `json:"user_id"`
}

func (favor *Favorite) CreateFavorite(db *gorm.DB) error {
	result := db.Create(&favor)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (favor *Favorite) DeleteFavorite(db *gorm.DB) error {
	result := db.Where("word_id = ? AND user_id = ?", favor.WordId, favor.UserId).Delete(&favor)
	if result.Error != nil {
		return result.Error
	}
	return nil
}