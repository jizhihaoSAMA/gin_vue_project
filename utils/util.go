package utils

import (
	"fmt"
	"gin_vue_project/dto"
	"math/rand"
	"time"
)

func RandomString(n int) string {
	var letters = []byte("asdfghjklzxcvbnmqwertyuiopASDFGHJKLZXCVBNMQWERTYUIOP")
	result := make([]byte, n)

	rand.Seed(time.Now().Unix())
	for i := range result {
		result[i] = letters[rand.Intn(len(letters))]
	}
	return string(result)
}

func InterfaceToUserDto(value interface{}) dto.UserDto {
	var u dto.UserDto
	switch v := value.(type) {
	case dto.UserDto:
		op, _ := value.(dto.UserDto)
		return op
	default:
		fmt.Println(v)
	}
	return u
}
