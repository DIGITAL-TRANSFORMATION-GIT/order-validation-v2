package models

type LoginForm struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Token struct {
	Token string `json:"token"`
}

type ChangePasswordForm struct {
	OldPassword string `json:"password"`
	NewPassword string `json:"new_password"`
}

type ChangeUsernameForm struct {
	NewUsername string `json:"new_username"`
	Password    string `json:"password"`
}
