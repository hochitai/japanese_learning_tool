package model

// ********************
// struct for swagger ui
type WordRequest struct {
	Characters    string `json:"characters" validate:"required"`
	Pronunciation string `json:"pronunciation" validate:"required"`
	Meaning       string `json:"meaning"`
}

type UserRequest struct {
	Name     string `json:"name"`
}

type UserInfoRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	Name     string `json:"name"`
}

type UserLoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type UserResponse struct {
	Id 	 int `json:"id"`
	Name string `json:"name"`
}

type UserLoginResponse struct {
	Name         string `json:"name"`
	AccessToken  string `json:"access-token"`
	RefreshToken string `json:"refresh-token"`
}

// *********************