package model

import "gorm.io/gorm"

type User struct {
	Id           int    `json:"id,omitempty"`
	Username     string `json:"username,omitempty"`
	Password     string `json:"password,omitempty"`
	Name         string `json:"name,omitempty"`
	Salt         string `json:"-"`
	AccessToken  string `json:"access-token" gorm:"-"`
	RefreshToken string `json:"refresh-token" gorm:"-"`
} 

func (user *User) GetUserByUsername(db *gorm.DB) (User, error) {
	var u User
	result := db.Where("username = ?", user.Username).First(&u)
	if result.Error != nil {
		return u, result.Error
	}
	return u, nil
}

func (user *User) CreateUser(db *gorm.DB) error {
	result := db.Create(&user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}