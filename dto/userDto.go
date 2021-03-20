package dto

import (
	"gin_vue_project/model"
	"strings"
)

type UserDto struct {
	ID        uint   `json:"id"`
	Username  string `json:"username"`
	Telephone string `json:"telephone"`
	Email     string `json:"email"`
	Detail    string `json:"detail"`
}

func ToUserDto(user model.User) UserDto {
	var emailDto string
	if user.Email != "" {
		emailDto = user.Email[:3] + "*****" + user.Email[strings.IndexByte(user.Email, '@'):]
	}

	return UserDto{
		ID:        user.ID,
		Username:  user.Username,
		Telephone: user.Telephone[:3] + "********",
		Email:     emailDto,
		Detail:    user.Detail,
	}
}
