package model

import "gorm.io/gorm"

type Favorite struct {
	Id     int `gorm:"primaryKey" json:"id"`
	WordId int `json:"word_id"`
	Word      Word `gorm:"foreignKey:WordId"`
	UserId int `json:"user_id"`
	User      User `gorm:"foreignKey:UserId"`
}

func (f *Favorite) GetFavorites(db *gorm.DB) ([]Word, error) {
	var words []Word
	result := db.Where("fav.user_id = ?", f.UserId).Joins("JOIN favorites fav on fav.word_id = words.id").Find(&words)
	if result.Error != nil {
		return words, result.Error
	}
	return words, nil
}

func (favor *Favorite) CreateFavorite(db *gorm.DB) error {
	result := db.Create(&favor)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// func (favor *Favorite) UpdateFavorite(

// )

func (favor *Favorite) DeleteFavorite(db *gorm.DB) error {
	result := db.Where("word_id = ? AND user_id = ?", favor.WordId, favor.UserId).Delete(&favor)
	if result.Error != nil {
		return result.Error
	}
	return nil
}