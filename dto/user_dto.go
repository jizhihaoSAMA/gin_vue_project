package dto

import "gin_vue_project/model"

type UserDto struct {
	ID        uint   `json:"id"`
	Username  string `json:"username"`
	Telephone string `json:"telephone"`
}

func ToUserDto(user model.User) UserDto {
	return UserDto{
		ID:        user.ID,
		Username:  user.Username,
		Telephone: user.Telephone[:3] + "********",
	}
}
