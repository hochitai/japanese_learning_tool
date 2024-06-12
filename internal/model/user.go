package model

import "gorm.io/gorm"

type User struct {
	Id           int    `gorm:"primaryKey" json:"id,omitempty"`
	Username     string `json:"username,omitempty"`
	Password     string `json:"password,omitempty"`
	Name         string `json:"name,omitempty"`
	Salt         string `json:"-"`
	Permission	 string `json:"permission,omitempty"` 
	AccessToken  string `json:"access-token,omitempty" gorm:"-"`
	RefreshToken string `json:"refresh-token,omitempty" gorm:"-"`
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
	user.Permission = "user"
	result := db.Create(&user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (user *User) CreateAdmin(db *gorm.DB) error {
	user.Permission = "admin"
	result := db.Create(&user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (user *User) UpdateUser(db *gorm.DB) error {
	// That time just only have 1 field of information can change
	// Will update more later
	result := db.Select("name").Updates(&user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (user *User) DeleteUser(db *gorm.DB) error {
	result := db.Delete(&user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (user User) GetUsers(db *gorm.DB) ([]User, error) {
	var users []User
	result := db.Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}