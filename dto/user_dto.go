package dto

import "gin_vue_project/model"

type UserDto struct {
	Username  string `json:"username"`
	Telephone string `json:"telephone"`
}

func ToUserDto(user model.User) UserDto {
	return UserDto{
		Username:  user.Username,
		Telephone: user.Telephone,
	}
}
